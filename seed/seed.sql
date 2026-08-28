INSERT INTO `seasonality` (`facility_id`,`seasonality_name`) VALUES
  (0x11F179CE2CF1844482674AA19E8DB8EA, 'Peak Season'),
  (0x11F179CE2CF1844482674AA19E8DB8EA, 'Off Season'),
  (0x11F179CE2CF1844482674AA19E8DB8EA, 'Shoulder Season'),
  (0x11F17A7142BA3475A40A9621B7EFE834, 'Peak Season'),
  (0x11F17A7142BA3475A40A9621B7EFE834, 'Off Season'),
  (0x11F17A7142BA3475A40A9621B7EFE834, 'Shoulder Season'),
  (0x11F17A715C23888AA40A9621B7EFE834, 'Peak Season'),
  (0x11F17A715C23888AA40A9621B7EFE834, 'Off Season'),
  (0x11F17A715C23888AA40A9621B7EFE834, 'Shoulder Season');


INSERT INTO `seasonality_number_of_people` (`seasonality_id`,`number_of_people`)
    WITH RECURSIVE seq AS (
        SELECT 1 AS n
        UNION ALL
        SELECT n + 1 FROM seq WHERE n < 10
    )
    SELECT s.id, seq.n
    FROM `seasonality` s
    CROSS JOIN seq
    WHERE s.seasonality_name IN ('Peak Season', 'Off Season', 'Shoulder Season');


UPDATE `seasonality_number_of_people` SET `price` = `number_of_people` * 10000.00


SELECT g.id AS customer_id,g.customer_name,g.customer_age,g.customer_sex,g.customer_from,middle.checkin FROM reservation 
INNER JOIN reservation_guest AS middle ON reservation.id = middle.reservation_id
INNER JOIN guest AS g ON middle.guest_id = g.id
WHERE reservation.facility_id = ?

INSERT INTO `seasonality_of_date` (`seasonality_id`,`date`,`room_type_id`)
    WITH RECURSIVE seq AS (
        SELECT '2026-08-01' AS n
        UNION ALL
        SELECT DATE_ADD(n,INTERVAL 1 DAY) FROM seq WHERE n < '2026-12-31'
    )
    SELECT s.id,seq.n,rt.id FROM `facilities` f 
    INNER JOIN `room` r ON f.id = r.facility_id
    INNER JOIN `room_type` rt ON r.id = rt.room_id
    INNER JOIN `seasonality` s ON f.id = s.facility_id
    CROSS JOIN seq ORDER BY seq.n


SELECT s.id,rt.id,s.id FROM `facilities` f 
INNER JOIN `room` r ON f.id = r.facility_id
INNER JOIN `room_type` rt ON r.id = rt.room_id
INNER JOIN `seasonality` s ON f.id = s.facility_id

SELECT r.room_name,rt.room_type_name,s.seasonality_name,snp.number_of_people,snp.price,sod.date FROM `facilities` f 
INNER JOIN `room` r ON f.id = r.facility_id
INNER JOIN `room_type` rt ON r.id = rt.room_id
INNER JOIN `seasonality` s ON f.id = s.facility_id
INNER JOIN `seasonality_number_of_people` snp ON s.id = snp.seasonality_id
INNER JOIN `seasonality_of_date` sod ON s.id = sod.seasonality_id
WHERE r.facility_id = ? AND room.id = ?;
