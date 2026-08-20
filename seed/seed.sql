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