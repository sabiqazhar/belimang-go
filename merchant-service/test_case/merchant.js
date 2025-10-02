const axios = require('axios');

// Set the correct port for the merchant service
const BASE_URL = 'http://localhost:8081';

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
            validateStatus: () => true, // Don't throw errors on non-2xx statuses
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
    log('🚀 Starting Merchant API Tests', 'blue');
    log('========================================\n', 'blue');

    const validPayload = {
        name: 'Warung Makan Sejahtera',
        merchantCategory: 'SmallRestaurant',
        imageUrl: 'https://example.com/image.png',
        location: {
            lat: -6.2088,
            long: 106.8456
        }
    };

    // --- SUCCESS CASE (201 CREATED) ---
    log('\n✅ SUCCESS CASE', 'blue');
    log('----------------------------------------\n', 'blue');

    await runTest(
        '1. POST /admin/merchants - Success',
        'POST',
        '/admin/merchants',
        validPayload,
        201
    );

    // --- VALIDATION FAILURE CASES (400 BAD REQUEST) ---
    log('\n❌ VALIDATION FAILURE CASES', 'blue');
    log('----------------------------------------\n', 'blue');

    // Name validation
    await runTest(
        '2. POST /admin/merchants - Name too short',
        'POST',
        '/admin/merchants',
        { ...validPayload, name: 'A' },
        400
    );
    await runTest(
        '3. POST /admin/merchants - Name too long',
        'POST',
        '/admin/merchants',
        { ...validPayload, name: 'This name is way too long and should fail validation for sure' },
        400
    );

    // MerchantCategory validation
    await runTest(
        '4. POST /admin/merchants - Invalid merchantCategory',
        'POST',
        '/admin/merchants',
        { ...validPayload, merchantCategory: 'InvalidCategory' },
        400
    );

    // ImageUrl validation
    await runTest(
        '5. POST /admin/merchants - Invalid imageUrl (not a URL)',
        'POST',
        '/admin/merchants',
        { ...validPayload, imageUrl: 'not-a-valid-url' },
        400
    );

    // Location validation
    await runTest(
        '6. POST /admin/merchants - Latitude is not a number',
        'POST',
        '/admin/merchants',
        { ...validPayload, location: { lat: 'invalid', long: 106.8456 } },
        400
    );
    await runTest(
        '7. POST /admin/merchants - Longitude is not a number',
        'POST',
        '/admin/merchants',
        { ...validPayload, location: { lat: -6.2088, long: 'invalid' } },
        400
    );

    // Missing fields validation
    const { name, ...payloadWithoutName } = validPayload;
    await runTest(
        '8. POST /admin/merchants - Missing name',
        'POST',
        '/admin/merchants',
        payloadWithoutName,
        400
    );

    const { merchantCategory, ...payloadWithoutCategory } = validPayload;
    await runTest(
        '9. POST /admin/merchants - Missing merchantCategory',
        'POST',
        '/admin/merchants',
        payloadWithoutCategory,
        400
    );

    const { imageUrl, ...payloadWithoutImageUrl } = validPayload;
    await runTest(
        '10. POST /admin/merchants - Missing imageUrl',
        'POST',
        '/admin/merchants',
        payloadWithoutImageUrl,
        400
    );

    const { location, ...payloadWithoutLocation } = validPayload;
    await runTest(
        '11. POST /admin/merchants - Missing location object',
        'POST',
        '/admin/merchants',
        payloadWithoutLocation,
        400
    );

    await runTest(
        '12. POST /admin/merchants - Missing latitude',
        'POST',
        '/admin/merchants',
        { ...validPayload, location: { long: 106.8456 } },
        400
    );

    await runTest(
        '13. POST /admin/merchants - Missing longitude',
        'POST',
        '/admin/merchants',
        { ...validPayload, location: { lat: -6.2088 } },
        400
    );

    await runTest(
        '14. POST /admin/merchants - Empty request body',
        'POST',
        '/admin/merchants',
        {},
        400
    );


    // --- SUMMARY ---
    log('\n========================================', 'blue');
    log('📊 TEST SUMMARY', 'blue');
    log('========================================', 'blue');
    log(`✓ Passed: ${passCount}`, 'green');
    log(`✗ Failed: ${failCount}`, 'red');
    log(`Total: ${passCount + failCount}`, 'blue');
    const successRate = passCount + failCount > 0 ? ((passCount / (passCount + failCount)) * 100).toFixed(2) : "0.00";
    log(`Success Rate: ${successRate}%`, 'yellow');
    log('========================================\n', 'blue');

    process.exit(failCount > 0 ? 1 : 0);
}

runAllTests();