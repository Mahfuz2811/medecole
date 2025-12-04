-- Migration: Add analytics counters to packages and exams tables
-- Date: 2025-07-31
-- Description: Add denormalized counter fields for better performance on enrollment and attempt statistics

-- Add analytics fields to packages table
ALTER TABLE packages 
ADD COLUMN enrollment_count INT DEFAULT 0 COMMENT 'Total number of users enrolled (includes free and paid)',
ADD COLUMN active_enrollment_count INT DEFAULT 0 COMMENT 'Current active enrollments (not expired)',
ADD COLUMN last_enrollment_at TIMESTAMP NULL COMMENT 'When the last user enrolled',
ADD INDEX idx_enrollment_count (enrollment_count);

-- Add analytics fields to exams table
ALTER TABLE exams
ADD COLUMN attempt_count INT DEFAULT 0 COMMENT 'Total number of exam attempts by all users',
ADD COLUMN completed_attempt_count INT DEFAULT 0 COMMENT 'Number of completed attempts (excludes abandoned)',
ADD COLUMN average_score DECIMAL(5,2) NULL COMMENT 'Average score of all completed attempts',
ADD COLUMN pass_rate DECIMAL(5,2) NULL COMMENT 'Percentage of attempts that passed',
ADD COLUMN last_attempt_at TIMESTAMP NULL COMMENT 'When the last attempt was made',
ADD INDEX idx_attempt_count (attempt_count);

-- Backfill existing data for packages
UPDATE packages p SET 
    enrollment_count = (
        SELECT COUNT(*) 
        FROM user_package_enrollment upe 
        WHERE upe.package_id = p.id
    ),
    active_enrollment_count = (
        SELECT COUNT(*) 
        FROM user_package_enrollment upe 
        WHERE upe.package_id = p.id 
        AND (upe.expires_at IS NULL OR upe.expires_at > NOW())
    ),
    last_enrollment_at = (
        SELECT MAX(upe.enrolled_at) 
        FROM user_package_enrollment upe 
        WHERE upe.package_id = p.id
    );

-- Backfill existing data for exams
UPDATE exams e SET 
    attempt_count = (
        SELECT COUNT(*) 
        FROM user_exam_attempt uea 
        WHERE uea.exam_id = e.id
    ),
    completed_attempt_count = (
        SELECT COUNT(*) 
        FROM user_exam_attempt uea 
        WHERE uea.exam_id = e.id 
        AND uea.status IN ('COMPLETED', 'AUTO_SUBMITTED')
        AND uea.is_scored = true
    ),
    average_score = (
        SELECT AVG(uea.score) 
        FROM user_exam_attempt uea 
        WHERE uea.exam_id = e.id 
        AND uea.status IN ('COMPLETED', 'AUTO_SUBMITTED')
        AND uea.is_scored = true
        AND uea.score IS NOT NULL
    ),
    pass_rate = (
        SELECT 
            CASE 
                WHEN COUNT(*) = 0 THEN NULL
                ELSE (COUNT(CASE WHEN uea.is_passed = true THEN 1 END) * 100.0 / COUNT(*))
            END
        FROM user_exam_attempt uea 
        WHERE uea.exam_id = e.id 
        AND uea.status IN ('COMPLETED', 'AUTO_SUBMITTED')
        AND uea.is_scored = true
    ),
    last_attempt_at = (
        SELECT MAX(uea.started_at) 
        FROM user_exam_attempt uea 
        WHERE uea.exam_id = e.id
    );

-- Create a stored procedure for updating package enrollment counts
DELIMITER //
CREATE PROCEDURE UpdatePackageEnrollmentCount(IN package_id_param INT, IN increment_value INT)
BEGIN
    UPDATE packages 
    SET 
        enrollment_count = enrollment_count + increment_value,
        last_enrollment_at = NOW()
    WHERE id = package_id_param;
END //
DELIMITER ;

-- Create a stored procedure for updating exam attempt counts  
DELIMITER //
CREATE PROCEDURE UpdateExamAttemptCount(IN exam_id_param INT)
BEGIN
    UPDATE exams 
    SET 
        attempt_count = attempt_count + 1,
        last_attempt_at = NOW()
    WHERE id = exam_id_param;
END //
DELIMITER ;

-- Create a stored procedure for updating exam completion stats
DELIMITER //
CREATE PROCEDURE UpdateExamCompletionStats(IN exam_id_param INT, IN score_param DECIMAL(5,2), IN passed_param BOOLEAN)
BEGIN
    DECLARE current_completed_count INT DEFAULT 0;
    DECLARE current_avg_score DECIMAL(5,2) DEFAULT 0;
    DECLARE current_pass_count INT DEFAULT 0;
    DECLARE new_pass_rate DECIMAL(5,2) DEFAULT 0;
    
    -- Increment completed attempt count
    UPDATE exams 
    SET completed_attempt_count = completed_attempt_count + 1
    WHERE id = exam_id_param;
    
    -- Recalculate average score and pass rate
    SELECT 
        COUNT(*) as completed_count,
        AVG(score) as avg_score,
        COUNT(CASE WHEN is_passed = true THEN 1 END) as pass_count
    INTO current_completed_count, current_avg_score, current_pass_count
    FROM user_exam_attempt 
    WHERE exam_id = exam_id_param 
    AND status IN ('COMPLETED', 'AUTO_SUBMITTED')
    AND is_scored = true;
    
    -- Calculate new pass rate
    IF current_completed_count > 0 THEN
        SET new_pass_rate = (current_pass_count * 100.0 / current_completed_count);
    END IF;
    
    -- Update exam statistics
    UPDATE exams 
    SET 
        average_score = current_avg_score,
        pass_rate = new_pass_rate
    WHERE id = exam_id_param;
END //
DELIMITER ;
