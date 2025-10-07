-- migrate:down

-- Delete from the 'child' table first to respect the foreign key constraint
DELETE FROM items WHERE id IN (
                               101, 102, 103, 104,
                               201, 202, 203, 204,
                               301, 302, 303, 304,
                               401, 402, 403,
                               501, 502, 503, 504,
                               601, 602, 603,
                               701, 702, 703,
                               801, 802, 803,
                               901, 902, 903,
                               1001, 1002, 1003,
                               1101, 1102, 1103,
                               1201, 1202, 1203,
                               1301, 1302, 1303,
                               1401, 1402, 1403,
                               1501, 1502, 1503
    );

-- Then delete from the 'parent' table
DELETE FROM merchants WHERE id IN (1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15);