const axios = require('axios');

const BASE_URL = 'http://localhost:8083';

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
            validateStatus: () => true,
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
    log('🚀 Starting Order Estimate API Tests', 'blue');
    log('========================================\n', 'blue');

    // ORDER ESTIMATE TESTS
    log('\n📋 ORDER ESTIMATE TESTS', 'blue');
    log('----------------------------------------\n', 'blue');

    // Valid request with single order
    await runTest(
        '1. POST /users/estimate - Success with single order',
        'POST',
        '/users/estimate',
        {
            userLocation: {
                lat: 1.0,
                long: 1.0
            },
            orders: [
                {
                    merchantId: "merchant-1",
                    isStartingPoint: true,
                    items: [
                        {
                            itemId: "item-1",
                            quantity: 2
                        }
                    ]
                }
            ]
        },
        200
    );

    // Valid request with multiple orders
    await runTest(
        '2. POST /users/estimate - Success with multiple orders',
        'POST',
        '/users/estimate',
        {
            userLocation: {
                lat: 1.5,
                long: 2.5
            },
            orders: [
                {
                    merchantId: "merchant-1",
                    isStartingPoint: true,
                    items: [
                        { itemId: "item-1", quantity: 1 },
                        { itemId: "item-2", quantity: 3 }
                    ]
                },
                {
                    merchantId: "merchant-2",
                    isStartingPoint: false,
                    items: [
                        { itemId: "item-3", quantity: 2 }
                    ]
                }
            ]
        },
        200
    );

    // Valid request with multiple items in single order
    await runTest(
        '3. POST /users/estimate - Success with multiple items',
        'POST',
        '/users/estimate',
        {
            userLocation: {
                lat: -6.2,
                long: 106.8
            },
            orders: [
                {
                    merchantId: "merchant-1",
                    isStartingPoint: true,
                    items: [
                        { itemId: "item-1", quantity: 1 },
                        { itemId: "item-2", quantity: 2 },
                        { itemId: "item-3", quantity: 3 }
                    ]
                }
            ]
        },
        200
    );

    // Missing userLocation
    await runTest(
        '4. POST /users/estimate - Missing userLocation',
        'POST',
        '/users/estimate',
        {
            orders: [
                {
                    merchantId: "merchant-1",
                    isStartingPoint: true,
                    items: [{ itemId: "item-1", quantity: 1 }]
                }
            ]
        },
        400
    );

    // Missing lat in userLocation
    await runTest(
        '5. POST /users/estimate - Missing lat',
        'POST',
        '/users/estimate',
        {
            userLocation: {
                long: 1.0
            },
            orders: [
                {
                    merchantId: "merchant-1",
                    isStartingPoint: true,
                    items: [{ itemId: "item-1", quantity: 1 }]
                }
            ]
        },
        400
    );

    // Missing long in userLocation
    await runTest(
        '6. POST /users/estimate - Missing long',
        'POST',
        '/users/estimate',
        {
            userLocation: {
                lat: 1.0
            },
            orders: [
                {
                    merchantId: "merchant-1",
                    isStartingPoint: true,
                    items: [{ itemId: "item-1", quantity: 1 }]
                }
            ]
        },
        400
    );

    // Missing orders array
    await runTest(
        '7. POST /users/estimate - Missing orders',
        'POST',
        '/users/estimate',
        {
            userLocation: {
                lat: 1.0,
                long: 1.0
            }
        },
        400
    );

    // Empty orders array
    await runTest(
        '8. POST /users/estimate - Empty orders array',
        'POST',
        '/users/estimate',
        {
            userLocation: {
                lat: 1.0,
                long: 1.0
            },
            orders: []
        },
        400
    );

    // Missing merchantId
    await runTest(
        '9. POST /users/estimate - Missing merchantId',
        'POST',
        '/users/estimate',
        {
            userLocation: {
                lat: 1.0,
                long: 1.0
            },
            orders: [
                {
                    isStartingPoint: true,
                    items: [{ itemId: "item-1", quantity: 1 }]
                }
            ]
        },
        400
    );

    // Missing isStartingPoint
    await runTest(
        '10. POST /users/estimate - Missing isStartingPoint',
        'POST',
        '/users/estimate',
        {
            userLocation: {
                lat: 1.0,
                long: 1.0
            },
            orders: [
                {
                    merchantId: "merchant-1",
                    items: [{ itemId: "item-1", quantity: 1 }]
                }
            ]
        },
        400
    );

    // Missing items array
    await runTest(
        '11. POST /users/estimate - Missing items',
        'POST',
        '/users/estimate',
        {
            userLocation: {
                lat: 1.0,
                long: 1.0
            },
            orders: [
                {
                    merchantId: "merchant-1",
                    isStartingPoint: true
                }
            ]
        },
        400
    );

    // Empty items array
    await runTest(
        '12. POST /users/estimate - Empty items array',
        'POST',
        '/users/estimate',
        {
            userLocation: {
                lat: 1.0,
                long: 1.0
            },
            orders: [
                {
                    merchantId: "merchant-1",
                    isStartingPoint: true,
                    items: []
                }
            ]
        },
        400
    );

    // Missing itemId in items
    await runTest(
        '13. POST /users/estimate - Missing itemId',
        'POST',
        '/users/estimate',
        {
            userLocation: {
                lat: 1.0,
                long: 1.0
            },
            orders: [
                {
                    merchantId: "merchant-1",
                    isStartingPoint: true,
                    items: [{ quantity: 1 }]
                }
            ]
        },
        400
    );

    // Missing quantity in items
    await runTest(
        '14. POST /users/estimate - Missing quantity',
        'POST',
        '/users/estimate',
        {
            userLocation: {
                lat: 1.0,
                long: 1.0
            },
            orders: [
                {
                    merchantId: "merchant-1",
                    isStartingPoint: true,
                    items: [{ itemId: "item-1" }]
                }
            ]
        },
        400
    );

    // No starting point (all false)
    await runTest(
        '15. POST /users/estimate - No starting point',
        'POST',
        '/users/estimate',
        {
            userLocation: {
                lat: 1.0,
                long: 1.0
            },
            orders: [
                {
                    merchantId: "merchant-1",
                    isStartingPoint: false,
                    items: [{ itemId: "item-1", quantity: 1 }]
                },
                {
                    merchantId: "merchant-2",
                    isStartingPoint: false,
                    items: [{ itemId: "item-2", quantity: 1 }]
                }
            ]
        },
        400
    );

    // Multiple starting points
    await runTest(
        '16. POST /users/estimate - Multiple starting points',
        'POST',
        '/users/estimate',
        {
            userLocation: {
                lat: 1.0,
                long: 1.0
            },
            orders: [
                {
                    merchantId: "merchant-1",
                    isStartingPoint: true,
                    items: [{ itemId: "item-1", quantity: 1 }]
                },
                {
                    merchantId: "merchant-2",
                    isStartingPoint: true,
                    items: [{ itemId: "item-2", quantity: 1 }]
                }
            ]
        },
        400
    );

    // Quantity zero
    await runTest(
        '17. POST /users/estimate - Quantity zero',
        'POST',
        '/users/estimate',
        {
            userLocation: {
                lat: 1.0,
                long: 1.0
            },
            orders: [
                {
                    merchantId: "merchant-1",
                    isStartingPoint: true,
                    items: [{ itemId: "item-1", quantity: 0 }]
                }
            ]
        },
        400
    );

    // Negative quantity
    await runTest(
        '18. POST /users/estimate - Negative quantity',
        'POST',
        '/users/estimate',
        {
            userLocation: {
                lat: 1.0,
                long: 1.0
            },
            orders: [
                {
                    merchantId: "merchant-1",
                    isStartingPoint: true,
                    items: [{ itemId: "item-1", quantity: -5 }]
                }
            ]
        },
        400
    );

    // Empty body
    await runTest(
        '19. POST /users/estimate - Empty body',
        'POST',
        '/users/estimate',
        {},
        400
    );

    // Invalid lat type (string)
    await runTest(
        '20. POST /users/estimate - Invalid lat type',
        'POST',
        '/users/estimate',
        {
            userLocation: {
                lat: "invalid",
                long: 1.0
            },
            orders: [
                {
                    merchantId: "merchant-1",
                    isStartingPoint: true,
                    items: [{ itemId: "item-1", quantity: 1 }]
                }
            ]
        },
        400
    );

    // Invalid long type (string)
    await runTest(
        '21. POST /users/estimate - Invalid long type',
        'POST',
        '/users/estimate',
        {
            userLocation: {
                lat: 1.0,
                long: "invalid"
            },
            orders: [
                {
                    merchantId: "merchant-1",
                    isStartingPoint: true,
                    items: [{ itemId: "item-1", quantity: 1 }]
                }
            ]
        },
        400
    );

    // Invalid quantity type (string)
    await runTest(
        '22. POST /users/estimate - Invalid quantity type',
        'POST',
        '/users/estimate',
        {
            userLocation: {
                lat: 1.0,
                long: 1.0
            },
            orders: [
                {
                    merchantId: "merchant-1",
                    isStartingPoint: true,
                    items: [{ itemId: "item-1", quantity: "invalid" }]
                }
            ]
        },
        400
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