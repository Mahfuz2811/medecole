# 🚀 Quizora Database Seeding Scripts

This directory contains scripts to populate your Quizora MCQ platform with comprehensive sample data.

## � Recent Updates

**The seeding system has been refactored for better maintainability!**

- See `README_REFACTORED.md` for details about the new modular structure
- Original monolithic seeder backed up as `seed_data_backup.go`
- New structure: Main orchestrator + separate seeders for users, packages, exams, and questions

## �📋 Available Scripts

### 1. **Go Seeder** (Recommended)

- **File**: `seed_data.go` (refactored modular version)
- **Runner**: `seed.sh`
- **Features**:
  - ✅ Comprehensive medical MCQ data
  - ✅ Realistic clinical scenarios
  - ✅ Proper data relationships
  - ✅ Automatic validation
  - ✅ Progress reporting
  - ✅ **NEW**: Modular, maintainable structure

### 2. **SQL Script** (Alternative)

- **File**: `sample_data.sql`
- **Features**:
  - ✅ Direct SQL insertion
  - ✅ Faster execution
  - ✅ Easy to customize
  - ✅ Cross-platform compatible

## 🎯 Sample Data Overview

### **4 Medical Subjects**

1. **Medicine** - Internal Medicine and specialties
2. **Surgery** - Surgical procedures and techniques
3. **Pediatrics** - Child and adolescent care
4. **Gynecology & Obstetrics** - Women's health

### **11 Body Systems**

- Cardiovascular System
- Respiratory System
- Gastrointestinal System
- Endocrine System
- General Surgery
- Orthopedic Surgery
- Neurosurgery
- Neonatology
- Pediatric Cardiology
- Obstetrics
- Gynecology

### **50+ MCQ Questions**

- **Question Types**: SBA (Single Best Answer) & True/False
- **Difficulty Levels**: 1 (Easy) to 3 (Hard)
- **Clinical Scenarios**: Realistic patient presentations
- **Detailed Explanations**: Educational content for each answer

## 🚀 How to Use

### Method 1: Go Seeder (Recommended)

```bash
# Navigate to backend directory
cd backend

# Run the seeding script
./scripts/seed.sh
```

### Method 2: SQL Script

```bash
# Connect to MySQL
mysql -u your_username -p

# Execute the SQL script
source scripts/sample_data.sql
```

### Method 3: Manual Go Execution

```bash
# From backend directory
cd backend

# Run the Go seeder directly
go run scripts/seed_data.go
```

## ⚙️ Prerequisites

1. **Database Setup**

   - MySQL server running
   - `quizora` database created
   - Proper credentials in config

2. **Application Setup**

   - Run `go run main.go` first to create tables
   - Ensure all migrations are complete

3. **Go Dependencies**
   - All Go modules installed (`go mod tidy`)

## 📊 Expected Results

After successful seeding:

- **4 subjects** with proper hierarchy
- **11 systems** categorized by subject
- **50+ questions** with realistic medical scenarios
- **JSON options** properly formatted
- **Foreign key relationships** intact

## 🔧 Troubleshooting

### Common Issues:

1. **Database Connection Failed**

   ```bash
   # Check MySQL is running
   brew services list | grep mysql

   # Check database exists
   mysql -u root -p -e "SHOW DATABASES LIKE 'quizora';"
   ```

2. **Table Doesn't Exist**

   ```bash
   # Run main application first
   go run main.go
   # Then run seeder
   ./scripts/seed.sh
   ```

3. **Permission Denied**

   ```bash
   # Make script executable
   chmod +x scripts/seed.sh
   ```

4. **Module Import Errors**
   ```bash
   # Install dependencies
   go mod tidy
   ```

## 🎨 Customization

### Adding More Questions

Edit `seed_data.go` and add questions to the respective system functions:

- `getCardiovascularQuestions()`
- `getRespiratoryQuestions()`
- `getGastrointestinalQuestions()`
- etc.

### Adding New Systems

1. Add system to subject in `getSampleSubjects()`
2. Create question generator function
3. Add questions with proper format

### Adding New Subjects

1. Add to `getSampleSubjects()` array
2. Create systems and questions
3. Update sort orders

## 📈 Data Structure

```
Subject (Medicine)
├── System (Cardiovascular)
│   ├── Question 1 (SBA)
│   ├── Question 2 (True/False)
│   └── Question 3 (SBA)
├── System (Respiratory)
│   ├── Question 1 (SBA)
│   └── Question 2 (True/False)
└── ...
```

## 🔍 Verification

After seeding, verify data:

```sql
-- Check record counts
SELECT COUNT(*) FROM subjects;
SELECT COUNT(*) FROM systems;
SELECT COUNT(*) FROM questions;

-- Check data distribution
SELECT
    s.name as subject_name,
    COUNT(sys.id) as systems_count,
    COUNT(q.id) as questions_count
FROM subjects s
LEFT JOIN systems sys ON s.id = sys.subject_id
LEFT JOIN questions q ON sys.id = q.system_id
GROUP BY s.id, s.name;
```

## 🌟 Next Steps

1. **Test Frontend**: Load frontend and browse questions
2. **API Testing**: Test question retrieval endpoints
3. **Quiz Features**: Implement quiz taking functionality
4. **Add More Data**: Expand with more medical specialties

---

**Happy Coding!** 🎯 Your MCQ platform is now loaded with realistic medical data.
