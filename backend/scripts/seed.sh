#!/bin/bash

# Data seeding script for Quizora MCQ Platform
# This script seeds the database with comprehensive sample data

echo "🚀 Starting Quizora Database Seeding..."
echo "======================================="

# Check if we're in the backend directory
if [ ! -f "go.mod" ]; then
    echo "❌ Error: Please run this script from the backend directory"
    exit 1
fi

# Check if the database is running
echo "📡 Checking database connection..."
go run -tags=scripts scripts/seed_data.go 2>/dev/null
if [ $? -eq 0 ]; then
    echo "✅ Database seeding completed successfully!"
    echo ""
    echo "📊 Sample data has been loaded:"
    echo "   • 4 Medical Subjects (Medicine, Surgery, Pediatrics, Gynecology & Obstetrics)"
    echo "   • 9 Body Systems (Cardiovascular, Respiratory, GI, Endocrine, etc.)"
    echo "   • 50+ MCQ Questions (SBA and True/False types)"
    echo ""
    echo "🎯 Your MCQ platform is now ready for testing!"
    echo ""
    echo "📝 Sample data includes:"
    echo "   • Realistic medical scenarios"
    echo "   • Detailed explanations"
    echo "   • Different difficulty levels"
    echo "   • Proper question categorization"
    echo ""
    echo "🌐 You can now:"
    echo "   • Test the frontend with real data"
    echo "   • Develop question browsing features"
    echo "   • Test quiz functionality"
    echo "   • Build reporting features"
else
    echo "❌ Database seeding failed!"
    echo ""
    echo "🔧 Troubleshooting steps:"
    echo "   1. Make sure MySQL is running"
    echo "   2. Check database credentials in config"
    echo "   3. Ensure the 'quizora' database exists"
    echo "   4. Run the main application first to create tables"
    echo ""
    echo "💡 Try running: go run main.go"
    echo "   Then run this script again"
    exit 1
fi

echo "✨ Happy coding! Your MCQ platform awaits..."
