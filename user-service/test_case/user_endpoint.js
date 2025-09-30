const axios = require('axios');

const BASE_URL = 'http://localhost:8080';
const ADMIN_USERNAME = 'admin12345';
const USER_USERNAME = 'user123456';
const ADMIN_EMAIL = 'admin@test.com';
const USER_EMAIL = 'user@test.com';

let passCount = 0;
let failCount = 0;

const colors = {
    reset: '\x1b[0m',
    green: '\x1b[32m',
    red: '\x1b[31m',
    yellow: '\x1b[33m',
    blue: '\x1b[34m',
};

function log(message, color = 'reset') {
    console.log(`${colors[color]}${message}${colors.reset}`);
}

async function runTest(testName, method, url, data, expectedStatus) {
    try {
        const response = await axios({
            method,
            url: `${BASE_URL}${url}`,
            data,
            validateStatus: () => true, // Don't throw on any status
        });

        const passed = response.status === expectedStatus;

        if (passed) {
            passCount++;
            log(`✓ ${testName}`, 'green');
            log(`  Expected: ${expectedStatus}, Got: ${response.status}`, 'green');
        } else {
            failCount++;
            log(`✗ ${testName}`, 'red');
            log(`  Expected: ${expectedStatus}, Got: ${response.status}`, 'red');
            if (response.data) {
                log(`  Response: ${JSON.stringify(response.data)}`, 'yellow');
            }
        }

        return { passed, response };
    } catch (error) {
        failCount++;
        log(`✗ ${testName}`, 'red');
        log(`  Error: ${error.message}`, 'red');
        return { passed: false, error };
    }
}

async function runAllTests() {
    log('\n========================================', 'blue');
    log('🚀 Starting API Tests', 'blue');
    log('========================================\n', 'blue');

    // ADMIN REGISTRATION TESTS
    log('\n📋 ADMIN REGISTRATION TESTS', 'blue');
    log('----------------------------------------\n', 'blue');

    await runTest(
        '1. POST /admin/register - Success',
        'POST',
        '/admin/register',
        { username: ADMIN_USERNAME, password: 'password123', email: ADMIN_EMAIL },
        201
    );

    await runTest(
        '2. POST /admin/register - Username too short',
        'POST',
        '/admin/register',
        { username: 'adm', password: 'password123', email: 'short@test.com' },
        400
    );

    await runTest(
        '3. POST /admin/register - Username too long',
        'POST',
        '/admin/register',
        { username: 'thisusernameiswaytoolongerthan30characters', password: 'password123', email: 'long@test.com' },
        400
    );

    await runTest(
        '4. POST /admin/register - Password too short',
        'POST',
        '/admin/register',
        { username: 'admin54321', password: 'pass', email: 'shortpass@test.com' },
        400
    );

    await runTest(
        '5. POST /admin/register - Password too long',
        'POST',
        '/admin/register',
        { username: 'admin11111', password: 'thispasswordiswaytoolongerthan30characters', email: 'longpass@test.com' },
        400
    );

    await runTest(
        '6. POST /admin/register - Invalid email format',
        'POST',
        '/admin/register',
        { username: 'admin99999', password: 'password123', email: 'invalidemail' },
        400
    );

    await runTest(
        '7. POST /admin/register - Missing username',
        'POST',
        '/admin/register',
        { password: 'password123', email: 'nousername@test.com' },
        400
    );

    await runTest(
        '8. POST /admin/register - Missing password',
        'POST',
        '/admin/register',
        { username: 'admin77777', email: 'nopassword@test.com' },
        400
    );

    await runTest(
        '9. POST /admin/register - Missing email',
        'POST',
        '/admin/register',
        { username: 'admin88888', password: 'password123' },
        400
    );

    await runTest(
        '10. POST /admin/register - Empty body',
        'POST',
        '/admin/register',
        {},
        400
    );

    await runTest(
        '11. POST /admin/register - Duplicate admin email',
        'POST',
        '/admin/register',
        { username: 'admin66666', password: 'password123', email: ADMIN_EMAIL },
        409
    );

    await runTest(
        '12. POST /admin/register - Duplicate username',
        'POST',
        '/admin/register',
        { username: ADMIN_USERNAME, password: 'password123', email: 'anotheradmin@test.com' },
        409
    );

    await runTest(
        '13. POST /admin/register - Username exactly 5 characters',
        'POST',
        '/admin/register',
        { username: 'adm55', password: 'password123', email: 'min5@test.com' },
        201
    );

    await runTest(
        '14. POST /admin/register - Username exactly 30 characters',
        'POST',
        '/admin/register',
        { username: 'adminthatisexactly30charlong', password: 'password123', email: 'max30@test.com' },
        201
    );

    await runTest(
        '15. POST /admin/register - Password exactly 5 characters',
        'POST',
        '/admin/register',
        { username: 'adminpass5', password: 'pass5', email: 'minpass@test.com' },
        201
    );

    await runTest(
        '16. POST /admin/register - Password exactly 30 characters',
        'POST',
        '/admin/register',
        { username: 'adminpass30', password: 'passwordthatisexactly30chars', email: 'maxpass@test.com' },
        201
    );

    // ADMIN LOGIN TESTS
    log('\n📋 ADMIN LOGIN TESTS', 'blue');
    log('----------------------------------------\n', 'blue');

    await runTest(
        '17. POST /admin/login - Success',
        'POST',
        '/admin/login',
        { username: ADMIN_USERNAME, password: 'password123' },
        200
    );

    await runTest(
        '18. POST /admin/login - Wrong password',
        'POST',
        '/admin/login',
        { username: ADMIN_USERNAME, password: 'wrongpassword' },
        400
    );

    await runTest(
        '19. POST /admin/login - Non-existent username',
        'POST',
        '/admin/login',
        { username: 'nonexistent12345', password: 'password123' },
        400
    );

    await runTest(
        '20. POST /admin/login - Username too short',
        'POST',
        '/admin/login',
        { username: 'adm', password: 'password123' },
        400
    );

    await runTest(
        '21. POST /admin/login - Password too short',
        'POST',
        '/admin/login',
        { username: ADMIN_USERNAME, password: 'pass' },
        400
    );

    await runTest(
        '22. POST /admin/login - Missing username',
        'POST',
        '/admin/login',
        { password: 'password123' },
        400
    );

    await runTest(
        '23. POST /admin/login - Missing password',
        'POST',
        '/admin/login',
        { username: ADMIN_USERNAME },
        400
    );

    await runTest(
        '24. POST /admin/login - Empty body',
        'POST',
        '/admin/login',
        {},
        400
    );

    // USER REGISTRATION TESTS
    log('\n📋 USER REGISTRATION TESTS', 'blue');
    log('----------------------------------------\n', 'blue');

    await runTest(
        '25. POST /users/register - Success',
        'POST',
        '/users/register',
        { username: USER_USERNAME, password: 'password123', email: USER_EMAIL },
        201
    );

    await runTest(
        '26. POST /users/register - Username too short',
        'POST',
        '/users/register',
        { username: 'usr', password: 'password123', email: 'short@test.com' },
        400
    );

    await runTest(
        '27. POST /users/register - Username too long',
        'POST',
        '/users/register',
        { username: 'thisusernameiswaytoolongerthan30characters', password: 'password123', email: 'long@test.com' },
        400
    );

    await runTest(
        '28. POST /users/register - Password too short',
        'POST',
        '/users/register',
        { username: 'user54321', password: 'pass', email: 'shortpass@test.com' },
        400
    );

    await runTest(
        '29. POST /users/register - Password too long',
        'POST',
        '/users/register',
        { username: 'user11111', password: 'thispasswordiswaytoolongerthan30characters', email: 'longpass@test.com' },
        400
    );

    await runTest(
        '30. POST /users/register - Invalid email format',
        'POST',
        '/users/register',
        { username: 'user99999', password: 'password123', email: 'invalidemail' },
        400
    );

    await runTest(
        '31. POST /users/register - Missing username',
        'POST',
        '/users/register',
        { password: 'password123', email: 'nousername@test.com' },
        400
    );

    await runTest(
        '32. POST /users/register - Missing password',
        'POST',
        '/users/register',
        { username: 'user77777', email: 'nopassword@test.com' },
        400
    );

    await runTest(
        '33. POST /users/register - Missing email',
        'POST',
        '/users/register',
        { username: 'user88888', password: 'password123' },
        400
    );

    await runTest(
        '34. POST /users/register - Empty body',
        'POST',
        '/users/register',
        {},
        400
    );

    await runTest(
        '35. POST /users/register - Duplicate user email',
        'POST',
        '/users/register',
        { username: 'user66666', password: 'password123', email: USER_EMAIL },
        409
    );

    await runTest(
        '36. POST /users/register - Duplicate username',
        'POST',
        '/users/register',
        { username: USER_USERNAME, password: 'password123', email: 'anotheruser@test.com' },
        409
    );

    await runTest(
        '37. POST /users/register - Username exactly 5 characters',
        'POST',
        '/users/register',
        { username: 'usr55', password: 'password123', email: 'usermin5@test.com' },
        201
    );

    await runTest(
        '38. POST /users/register - Username exactly 30 characters',
        'POST',
        '/users/register',
        { username: 'usernamethatisexactly30charl', password: 'password123', email: 'usermax30@test.com' },
        201
    );

    await runTest(
        '39. POST /users/register - Password exactly 5 characters',
        'POST',
        '/users/register',
        { username: 'userpass5', password: 'pass5', email: 'userminpass@test.com' },
        201
    );

    await runTest(
        '40. POST /users/register - Password exactly 30 characters',
        'POST',
        '/users/register',
        { username: 'userpass30', password: 'passwordthatisexactly30chars', email: 'usermaxpass@test.com' },
        201
    );

    // USER LOGIN TESTS
    log('\n📋 USER LOGIN TESTS', 'blue');
    log('----------------------------------------\n', 'blue');

    await runTest(
        '41. POST /users/login - Success',
        'POST',
        '/users/login',
        { username: USER_USERNAME, password: 'password123' },
        200
    );

    await runTest(
        '42. POST /users/login - Wrong password',
        'POST',
        '/users/login',
        { username: USER_USERNAME, password: 'wrongpassword' },
        400
    );

    await runTest(
        '43. POST /users/login - Non-existent username',
        'POST',
        '/users/login',
        { username: 'nonexistent12345', password: 'password123' },
        400
    );

    await runTest(
        '44. POST /users/login - Username too short',
        'POST',
        '/users/login',
        { username: 'usr', password: 'password123' },
        400
    );

    await runTest(
        '45. POST /users/login - Password too short',
        'POST',
        '/users/login',
        { username: USER_USERNAME, password: 'pass' },
        400
    );

    await runTest(
        '46. POST /users/login - Missing username',
        'POST',
        '/users/login',
        { password: 'password123' },
        400
    );

    await runTest(
        '47. POST /users/login - Missing password',
        'POST',
        '/users/login',
        { username: USER_USERNAME },
        400
    );

    await runTest(
        '48. POST /users/login - Empty body',
        'POST',
        '/users/login',
        {},
        400
    );

    // CROSS-TYPE TESTS
    log('\n📋 CROSS-TYPE USERNAME CONFLICT TESTS', 'blue');
    log('----------------------------------------\n', 'blue');

    await runTest(
        '49. POST /users/register - Username conflict with admin',
        'POST',
        '/users/register',
        { username: ADMIN_USERNAME, password: 'password123', email: 'usertryingadminname@test.com' },
        409
    );

    await runTest(
        '50. POST /admin/register - Username conflict with user',
        'POST',
        '/admin/register',
        { username: USER_USERNAME, password: 'password123', email: 'admintryingusername@test.com' },
        409
    );

    // EMAIL CROSS-TYPE TESTS
    log('\n📋 EMAIL CROSS-TYPE TESTS (Should succeed)', 'blue');
    log('----------------------------------------\n', 'blue');

    await runTest(
        '51. POST /admin/register - Same email as user',
        'POST',
        '/admin/register',
        { username: 'adminsameemail', password: 'password123', email: USER_EMAIL },
        201
    );

    await runTest(
        '52. POST /users/register - Same email as admin',
        'POST',
        '/users/register',
        { username: 'usersameemail', password: 'password123', email: ADMIN_EMAIL },
        201
    );

    // SUMMARY
    log('\n========================================', 'blue');
    log('📊 TEST SUMMARY', 'blue');
    log('========================================', 'blue');
    log(`✓ Passed: ${passCount}`, 'green');
    log(`✗ Failed: ${failCount}`, 'red');
    log(`Total: ${passCount + failCount}`, 'blue');
    log(`Success Rate: ${((passCount / (passCount + failCount)) * 100).toFixed(2)}%`, 'yellow');
    log('========================================\n', 'blue');

    process.exit(failCount > 0 ? 1 : 0);
}

runAllTests();