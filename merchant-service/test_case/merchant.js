const axios = require('axios');

// Set the correct port for the merchant service
const BASE_URL = 'http://localhost:8081';

let passCount = 0;
let failCount = 0;
let createdMerchants = [];

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

async function runTest(testName, method, url, data, expectedStatus, validateResponse = null) {
    try {
        const response = await axios({
            method,
            url: `${BASE_URL}${url}`,
            data,
            validateStatus: () => true,
        });

        let passed = response.status === expectedStatus;

        // Additional response validation if provided
        if (passed && validateResponse) {
            const validationResult = validateResponse(response.data);
            if (!validationResult.valid) {
                passed = false;
                log(`✗ ${testName}`, 'red');
                log(`  Status: ${response.status} ✓`, 'green');
                log(`  Validation: ${validationResult.message}`, 'red');
                failCount++;
                return { passed, response };
            }
        }

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

async function createTestMerchant(name, category = 'SmallRestaurant') {
    const payload = {
        name,
        merchantCategory: category,
        imageUrl: 'https://example.com/image.png',
        location: {
            lat: -6.2088,
            long: 106.8456
        }
    };

    try {
        const response = await axios.post(`${BASE_URL}/admin/merchants`, payload);
        if (response.status === 201 && response.data.merchantId) {
            createdMerchants.push({
                merchantId: response.data.merchantId,
                name,
                category
            });
            return response.data.merchantId;
        }
    } catch (error) {
        log(`Failed to create test merchant: ${name}`, 'red');
    }
    return null;
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
    log('\n✅ POST SUCCESS CASE', 'blue');
    log('----------------------------------------\n', 'blue');

    await runTest(
        '1. POST /admin/merchants - Success',
        'POST',
        '/admin/merchants',
        validPayload,
        201
    );

    // --- VALIDATION FAILURE CASES (400 BAD REQUEST) ---
    log('\n❌ POST VALIDATION FAILURE CASES', 'blue');
    log('----------------------------------------\n', 'blue');

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

    await runTest(
        '4. POST /admin/merchants - Invalid merchantCategory',
        'POST',
        '/admin/merchants',
        { ...validPayload, merchantCategory: 'InvalidCategory' },
        400
    );

    await runTest(
        '5. POST /admin/merchants - Invalid imageUrl (not a URL)',
        'POST',
        '/admin/merchants',
        { ...validPayload, imageUrl: 'not-a-valid-url' },
        400
    );

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

    // --- CREATE TEST DATA FOR GET MERCHANTS ---
    log('\n🔧 Creating Test Merchants for GET Tests', 'blue');
    log('----------------------------------------\n', 'blue');

    await createTestMerchant('Kayleen Restaurant', 'SmallRestaurant');
    await createTestMerchant('Excellent Cafe', 'MediumRestaurant');
    await createTestMerchant('Green Valley Store', 'LargeRestaurant');
    await createTestMerchant('Booth Paradise', 'BoothKiosk');
    await createTestMerchant('Merchandise Hub', 'MerchandiseRestaurant');
    await createTestMerchant('Quick Stop', 'ConvenienceStore');

    log(`Created ${createdMerchants.length} test merchants\n`, 'green');

    // --- GET MERCHANTS TESTS ---
    log('\n✅ GET MERCHANTS - SUCCESS CASES', 'blue');
    log('----------------------------------------\n', 'blue');

    const result15 = await runTest(
        '15. GET /admin/merchants - Get all merchants (default pagination)',
        'GET',
        '/admin/merchants',
        null,
        200,
        (response) => {
            if (!response.data || !Array.isArray(response.data)) {
                return { valid: false, message: 'Response should have data array' };
            }
            if (!response.meta || typeof response.meta !== 'object') {
                return { valid: false, message: 'Response should have meta object' };
            }
            if (typeof response.meta.limit !== 'number' ||
                typeof response.meta.offset !== 'number' ||
                typeof response.meta.total !== 'number') {
                return { valid: false, message: 'Meta should have limit, offset, and total as numbers' };
            }
            if (response.meta.limit !== 5 || response.meta.offset !== 0) {
                return { valid: false, message: 'Default should be limit=5, offset=0' };
            }
            return { valid: true };
        }
    );
    if (result15.response?.data) {
        log(`  Found ${result15.response.data.data?.length || 0} merchants (total: ${result15.response.data.meta?.total || 0})`, 'yellow');
    }

    const result16 = await runTest(
        '16. GET /admin/merchants - With limit and offset',
        'GET',
        '/admin/merchants?limit=3&offset=0',
        null,
        200,
        (response) => {
            if (!response.data || !Array.isArray(response.data)) {
                return { valid: false, message: 'Response should have data array' };
            }
            if (!response.meta || typeof response.meta !== 'object') {
                return { valid: false, message: 'Response should have meta object' };
            }
            if (response.data.length > 3) {
                return { valid: false, message: `Expected max 3 merchants, got ${response.data.length}` };
            }
            if (response.meta.limit !== 3 || response.meta.offset !== 0) {
                return { valid: false, message: 'Meta should reflect limit=3, offset=0' };
            }
            return { valid: true };
        }
    );
    if (result16.response?.data) {
        log(`  Returned ${result16.response.data.data?.length || 0} merchants (limit=3)`, 'yellow');
    }

    await runTest(
        '18. GET /admin/merchants - Filter by merchantId (non-existent)',
        'GET',
        '/admin/merchants?merchantId=999999999',
        null,
        200,
        (response) => {
            if (!response.data || !Array.isArray(response.data)) {
                return { valid: false, message: 'Response should have data array' };
            }
            if (!response.meta) {
                return { valid: false, message: 'Response should have meta object' };
            }
            if (response.data.length !== 0) {
                return { valid: false, message: 'Should return empty array for non-existent ID' };
            }
            if (response.meta.total !== 0) {
                return { valid: false, message: 'Total should be 0 for non-existent ID' };
            }
            return { valid: true };
        }
    );

    await runTest(
        '19. GET /admin/merchants - Filter by name (wildcard search "een")',
        'GET',
        '/admin/merchants?name=een',
        null,
        200,
        (response) => {
            if (!response.data || !Array.isArray(response.data)) {
                return { valid: false, message: 'Response should have data array' };
            }
            if (!response.meta) {
                return { valid: false, message: 'Response should have meta object' };
            }
            const hasMatchingName = response.data.every(m =>
                m.name && m.name.toLowerCase().includes('een')
            );
            if (response.data.length > 0 && !hasMatchingName) {
                return { valid: false, message: 'All returned merchants should have "een" in name' };
            }
            return { valid: true };
        }
    );

    await runTest(
        '20. GET /admin/merchants - Filter by name (case insensitive)',
        'GET',
        '/admin/merchants?name=KAYLEEN',
        null,
        200,
        (response) => {
            if (!response.data || !Array.isArray(response.data)) {
                return { valid: false, message: 'Response should have data array' };
            }
            if (!response.meta) {
                return { valid: false, message: 'Response should have meta object' };
            }
            const hasMatchingName = response.data.every(m =>
                m.name && m.name.toLowerCase().includes('kayleen')
            );
            if (response.data.length > 0 && !hasMatchingName) {
                return { valid: false, message: 'Search should be case insensitive' };
            }
            return { valid: true };
        }
    );

    await runTest(
        '21. GET /admin/merchants - Filter by name (non-existent)',
        'GET',
        '/admin/merchants?name=NonExistentMerchantXYZ',
        null,
        200,
        (response) => {
            if (!response.data || !Array.isArray(response.data)) {
                return { valid: false, message: 'Response should have data array' };
            }
            if (!response.meta) {
                return { valid: false, message: 'Response should have meta object' };
            }
            if (response.data.length !== 0) {
                return { valid: false, message: 'Should return empty array for non-existent name' };
            }
            if (response.meta.total !== 0) {
                return { valid: false, message: 'Total should be 0 for non-existent name' };
            }
            return { valid: true };
        }
    );

    await runTest(
        '22. GET /admin/merchants - Filter by merchantCategory (SmallRestaurant)',
        'GET',
        '/admin/merchants?merchantCategory=SmallRestaurant',
        null,
        200,
        (response) => {
            if (!response.data || !Array.isArray(response.data)) {
                return { valid: false, message: 'Response should have data array' };
            }
            if (!response.meta) {
                return { valid: false, message: 'Response should have meta object' };
            }
            const allCorrectCategory = response.data.every(m =>
                m.merchantCategory === 'SmallRestaurant'
            );
            if (response.data.length > 0 && !allCorrectCategory) {
                return { valid: false, message: 'All returned merchants should be SmallRestaurant' };
            }
            return { valid: true };
        }
    );

    await runTest(
        '23. GET /admin/merchants - Filter by merchantCategory (BoothKiosk)',
        'GET',
        '/admin/merchants?merchantCategory=BoothKiosk',
        null,
        200,
        (response) => {
            if (!response.data || !Array.isArray(response.data)) {
                return { valid: false, message: 'Response should have data array' };
            }
            if (!response.meta) {
                return { valid: false, message: 'Response should have meta object' };
            }
            const allCorrectCategory = response.data.every(m =>
                m.merchantCategory === 'BoothKiosk'
            );
            if (response.data.length > 0 && !allCorrectCategory) {
                return { valid: false, message: 'All returned merchants should be BoothKiosk' };
            }
            return { valid: true };
        }
    );

    await runTest(
        '24. GET /admin/merchants - Filter by merchantCategory (invalid enum - should return empty)',
        'GET',
        '/admin/merchants?merchantCategory=InvalidCategory',
        null,
        200,
        (response) => {
            if (!response.data || !Array.isArray(response.data)) {
                return { valid: false, message: 'Response should have data array' };
            }
            if (!response.meta) {
                return { valid: false, message: 'Response should have meta object' };
            }
            if (response.data.length !== 0) {
                return { valid: false, message: 'Should return empty array for invalid category' };
            }
            if (response.meta.total !== 0) {
                return { valid: false, message: 'Total should be 0 for invalid category' };
            }
            return { valid: true };
        }
    );

    await runTest(
        '25. GET /admin/merchants - Sort by createdAt (asc)',
        'GET',
        '/admin/merchants?createdAt=asc&limit=10',
        null,
        200,
        (response) => {
            if (!response.data || !Array.isArray(response.data)) {
                return { valid: false, message: 'Response should have data array' };
            }
            if (!response.meta) {
                return { valid: false, message: 'Response should have meta object' };
            }
            if (response.data.length > 1) {
                for (let i = 1; i < response.data.length; i++) {
                    const prev = new Date(response.data[i - 1].createdAt);
                    const curr = new Date(response.data[i].createdAt);
                    if (prev > curr) {
                        return { valid: false, message: 'Merchants should be sorted by createdAt ascending' };
                    }
                }
            }
            return { valid: true };
        }
    );

    await runTest(
        '26. GET /admin/merchants - Sort by createdAt (desc)',
        'GET',
        '/admin/merchants?createdAt=desc&limit=10',
        null,
        200,
        (response) => {
            if (!response.data || !Array.isArray(response.data)) {
                return { valid: false, message: 'Response should have data array' };
            }
            if (!response.meta) {
                return { valid: false, message: 'Response should have meta object' };
            }
            if (response.data.length > 1) {
                for (let i = 1; i < response.data.length; i++) {
                    const prev = new Date(response.data[i - 1].createdAt);
                    const curr = new Date(response.data[i].createdAt);
                    if (prev < curr) {
                        return { valid: false, message: 'Merchants should be sorted by createdAt descending' };
                    }
                }
            }
            return { valid: true };
        }
    );

    await runTest(
        '27. GET /admin/merchants - Sort by createdAt (invalid value - should ignore)',
        'GET',
        '/admin/merchants?createdAt=invalid&limit=5',
        null,
        200,
        (response) => {
            if (!response.data || !Array.isArray(response.data)) {
                return { valid: false, message: 'Response should have data array' };
            }
            if (!response.meta) {
                return { valid: false, message: 'Response should have meta object' };
            }
            return { valid: true };
        }
    );

    await runTest(
        '28. GET /admin/merchants - Multiple filters (name + category)',
        'GET',
        '/admin/merchants?name=restaurant&merchantCategory=SmallRestaurant',
        null,
        200,
        (response) => {
            if (!response.data || !Array.isArray(response.data)) {
                return { valid: false, message: 'Response should have data array' };
            }
            if (!response.meta) {
                return { valid: false, message: 'Response should have meta object' };
            }
            const allMatch = response.data.every(m =>
                m.name && m.name.toLowerCase().includes('restaurant') &&
                m.merchantCategory === 'SmallRestaurant'
            );
            if (response.data.length > 0 && !allMatch) {
                return { valid: false, message: 'All merchants should match both filters' };
            }
            return { valid: true };
        }
    );

    await runTest(
        '29. GET /admin/merchants - Invalid limit (non-numeric)',
        'GET',
        '/admin/merchants?limit=abc',
        null,
        400
    );

    await runTest(
        '30. GET /admin/merchants - Invalid offset (non-numeric)',
        'GET',
        '/admin/merchants?offset=xyz',
        null,
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