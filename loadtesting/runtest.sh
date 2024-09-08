#todo refresh db
sqlite3 ../test.db ".read ../migrations/0-init.sql"
k6 run --out json=test_results.json hourTest.js