-- migrate:up

-- ----------------------------
-- DUMMY DATA FOR MERCHANTS (Corrected Categories)
-- ----------------------------
INSERT INTO merchants (id, name, merchant_category, image_url, longitude, latitude, h3_index) VALUES
-- Mapping: CoffeeShops/BubbleTea -> BoothKiosk, Restaurants -> Respective Size, Groceries -> ConvenienceStore, Retail -> MerchandiseRestaurant
(1, 'Kopi Kenangan - Alsut', 'BoothKiosk', 'https://example.com/images/kopi_kenangan.jpg', 106.6534, -6.2244, 892834537492840447),
(2, 'Warung Nasi Ibu Siti', 'SmallRestaurant', 'https://example.com/images/warung_siti.jpg', 106.6601, -6.2291, 892834537492840448),
(3, 'Indomaret Point Gading Serpong', 'ConvenienceStore', 'https://example.com/images/indomaret.jpg', 106.6288, -6.2415, 892834537492840455),
(4, 'Bakerman - Central Park', 'SmallRestaurant', 'https://example.com/images/bakerman.jpg', 106.7915, -6.1766, 892834537492840449),
(5, 'Sate Khas Senayan - PIM', 'LargeRestaurant', 'https://example.com/images/sate_senayan.jpg', 106.7997, -6.2447, 892834537492840450),
(6, 'McDonald''s - Bintaro', 'MediumRestaurant', 'https://example.com/images/mcd.jpg', 106.7325, -6.2808, 892834537492840451),
(7, 'Chatime - Living World', 'BoothKiosk', 'https://example.com/images/chatime.jpg', 106.6578, -6.2498, 892834537492840452),
(8, 'Geprek Bensu - Karawaci', 'SmallRestaurant', 'https://example.com/images/geprek_bensu.jpg', 106.6212, -6.2309, 892834537492840453),
(9, 'Gramedia - Grand Indonesia', 'MerchandiseRestaurant', 'https://example.com/images/gramedia.jpg', 106.8220, -6.1953, 892834537492840456),
(10, 'Guardian Pharmacy - CP', 'MerchandiseRestaurant', 'https://example.com/images/guardian.jpg', 106.7910, -6.1772, 892834537492840457),
(11, 'Uniqlo - PIK Avenue', 'MerchandiseRestaurant', 'https://example.com/images/uniqlo.jpg', 106.7418, -6.1098, 892834537492840458),
(12, 'Toko Kelontong Pak Budi', 'ConvenienceStore', 'https://example.com/images/kelontong.jpg', 106.7011, -6.2555, 892834537492840459),
(13, 'ACE Hardware - Gandaria City', 'MerchandiseRestaurant', 'https://example.com/images/ace.jpg', 106.7830, -6.2443, 892834537492840460),
(14, 'Digimap - Senayan City', 'MerchandiseRestaurant', 'https://example.com/images/digimap.jpg', 106.7981, -6.2273, 892834537492840461),
(15, 'Starbucks - Summarecon Mall', 'BoothKiosk', 'https://example.com/images/starbucks.jpg', 106.6280, -6.2429, 892834537492840454);

-- ----------------------------
-- DUMMY DATA FOR ITEMS (No change needed here, as product categories are distinct from merchant categories)
-- ----------------------------
INSERT INTO items (id, merchant_id, name, product_category, image_url, price) VALUES
-- Items for Kopi Kenangan (merchant_id: 1)
(101, 1, 'Kopi Kenangan Mantan', 'Beverage', 'https://example.com/items/kopi_mantan.jpg', 18000.00),
(102, 1, 'Americano', 'Beverage', 'https://example.com/items/americano.jpg', 15000.00),
(103, 1, 'Donat Gula', 'Snack', 'https://example.com/items/donat.jpg', 9000.00),
(104, 1, 'Teh Susu Gula Aren', 'Beverage', 'https://example.com/items/teh_susu.jpg', 19000.00),

-- Items for Warung Nasi Ibu Siti (merchant_id: 2)
(201, 2, 'Nasi Ayam Goreng', 'MainCourse', 'https://example.com/items/nasi_ayam.jpg', 25000.00),
(202, 2, 'Mie Goreng Spesial', 'MainCourse', 'https://example.com/items/mie_goreng.jpg', 22000.00),
(203, 2, 'Es Teh Manis', 'Beverage', 'https://example.com/items/es_teh.jpg', 5000.00),
(204, 2, 'Tahu Goreng', 'Snack', 'https://example.com/items/tahu.jpg', 8000.00),

-- Items for Indomaret (merchant_id: 3)
(301, 3, 'Indomie Goreng', 'InstantFood', 'https://example.com/items/indomie.jpg', 3500.00),
(302, 3, 'Aqua 600ml', 'Beverage', 'https://example.com/items/aqua.jpg', 4000.00),
(303, 3, 'Chitato Sapi Panggang', 'Snack', 'https://example.com/items/chitato.jpg', 11000.00),
(304, 3, 'Sari Roti Coklat', 'Bakery', 'https://example.com/items/sari_roti.jpg', 6000.00),

-- Items for Bakerman (merchant_id: 4)
(401, 4, 'Plain Croissant', 'Pastry', 'https://example.com/items/croissant.jpg', 20000.00),
(402, 4, 'Cinnamon Roll', 'Pastry', 'https://example.com/items/cinnamon_roll.jpg', 28000.00),
(403, 4, 'Kouign Amann', 'Pastry', 'https://example.com/items/kouign_amann.jpg', 32000.00),

-- Items for Sate Khas Senayan (merchant_id: 5)
(501, 5, 'Sate Ayam (10 tusuk)', 'MainCourse', 'https://example.com/items/sate_ayam.jpg', 75000.00),
(502, 5, 'Lontong', 'SideDish', 'https://example.com/items/lontong.jpg', 10000.00),
(503, 5, 'Tahu Telor', 'MainCourse', 'https://example.com/items/tahu_telor.jpg', 55000.00),
(504, 5, 'Soto Betawi', 'MainCourse', 'https://example.com/items/soto_betawi.jpg', 68000.00),

-- Items for McDonald's (merchant_id: 6)
(601, 6, 'Paket Panas 1', 'ComboMeal', 'https://example.com/items/panas1.jpg', 35000.00),
(602, 6, 'Big Mac', 'Burger', 'https://example.com/items/big_mac.jpg', 42000.00),
(603, 6, 'McFlurry Oreo', 'Dessert', 'https://example.com/items/mcflurry.jpg', 14000.00),

-- Items for Chatime (merchant_id: 7)
(701, 7, 'Pearl Milk Tea', 'BubbleTea', 'https://example.com/items/pearl_milk_tea.jpg', 25000.00),
(702, 7, 'Hazelnut Chocolate Milk Tea', 'BubbleTea', 'https://example.com/items/hazelnut.jpg', 28000.00),
(703, 7, 'Grass Jelly Roasted Milk Tea', 'BubbleTea', 'https://example.com/items/grass_jelly.jpg', 26000.00),

-- Items for Geprek Bensu (merchant_id: 8)
(801, 8, 'Paket Geprek Nasi', 'MainCourse', 'https://example.com/items/geprek_nasi.jpg', 20000.00),
(802, 8, 'Paket Geprek Indomie', 'MainCourse', 'https://example.com/items/geprek_indomie.jpg', 22000.00),
(803, 8, 'Kulit Ayam Crispy', 'Snack', 'https://example.com/items/kulit_ayam.jpg', 15000.00),

-- Items for Gramedia (merchant_id: 9)
(901, 9, 'Novel "Laskar Pelangi"', 'Book', 'https://example.com/items/laskar_pelangi.jpg', 85000.00),
(902, 9, 'Buku Tulis Sidu 38 Lembar', 'Stationery', 'https://example.com/items/sidu.jpg', 4000.00),
(903, 9, 'Pulpen Pilot G2', 'Stationery', 'https://example.com/items/pilot_g2.jpg', 21000.00),

-- Items for Guardian (merchant_id: 10)
(1001, 10, 'Panadol Biru', 'Medicine', 'https://example.com/items/panadol.jpg', 11500.00),
(1002, 10, 'Hansaplast Plester', 'Health', 'https://example.com/items/hansaplast.jpg', 8000.00),
(1003, 10, 'Vitamin C IPI', 'Supplement', 'https://example.com/items/vit_c.jpg', 7500.00),

-- Items for Uniqlo (merchant_id: 11)
(1101, 11, 'UT T-Shirt', 'Apparel', 'https://example.com/items/ut.jpg', 249000.00),
(1102, 11, 'Celana Jeans Slim Fit', 'Apparel', 'https://example.com/items/jeans.jpg', 599000.00),
(1103, 11, 'Jaket Hoodie', 'Apparel', 'https://example.com/items/hoodie.jpg', 499000.00),

-- Items for Toko Pak Budi (merchant_id: 12)
(1201, 12, 'Telur Ayam 1kg', 'Groceries', 'https://example.com/items/telur.jpg', 28000.00),
(1202, 12, 'Gas Elpiji 3kg', 'Household', 'https://example.com/items/gas.jpg', 22000.00),
(1203, 12, 'Beras 5kg', 'Groceries', 'https://example.com/items/beras.jpg', 65000.00),

-- Items for ACE Hardware (merchant_id: 13)
(1301, 13, 'Set Obeng Krisbow', 'Tools', 'https://example.com/items/obeng.jpg', 150000.00),
(1302, 13, 'Lampu Bohlam LED', 'Household', 'https://example.com/items/bohlam.jpg', 45000.00),
(1303, 13, 'Rak Penyimpanan 4 Tingkat', 'Furniture', 'https://example.com/items/rak.jpg', 450000.00),

-- Items for Digimap (merchant_id: 14)
(1401, 14, 'iPhone 15 Pro', 'Electronics', 'https://example.com/items/iphone15.jpg', 18999000.00),
(1402, 14, 'AirPods Pro 2', 'Accessories', 'https://example.com/items/airpods.jpg', 3999000.00),
(1403, 14, 'Kabel USB-C', 'Accessories', 'https://example.com/items/usbc.jpg', 349000.00),

-- Items for Starbucks (merchant_id: 15)
(1501, 15, 'Caffe Latte', 'Beverage', 'https://example.com/items/latte.jpg', 48000.00),
(1502, 15, 'Caramel Macchiato', 'Beverage', 'https://example.com/items/macchiato.jpg', 59000.00),
(1503, 15, 'Classic Croissant', 'Pastry', 'https://example.com/items/sbux_croissant.jpg', 25000.00);