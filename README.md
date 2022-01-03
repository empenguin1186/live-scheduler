# データベース設計

## Live テーブル
- id: ライブのID(Auto Increment, 主キー)
- name: ライブ名
- location: ライブの開催場所
- date: ライブ開催日
- performance_fee: 一人当たりの出演費
- equipment_cost: １バンドあたりの機材費

## Band テーブル
- name: バンド名
- live_id: ライブID(主キー) 外部キーとして Live テーブルの id カラムを参照する
- turn: 出演順(主キー)

## BandMember テーブル
- live_id: ライブID(主キー) Live テーブルの id カラムを外部キー
- turn: 出演順(主キー) Band テーブルの order カラムを外部キー
- member_name: メンバーの名前(主キー) Member テーブルの name カラムを外部キー
- member_part: 担当パート(主キー) Member テーブルの part カラムを外部キー

## Player テーブル
- name: メンバーの名前(主キー)
- part: 担当パート(主キー) enum 型

```mysql
# テーブル作成
CREATE TABLE Live ( id SERIAL PRIMARY KEY, name VARCHAR(50), location VARCHAR(50), date DATE, performance_fee INT, equipment_cost INT );
CREATE TABLE Band ( name VARCHAR(50), live_id BIGINT UNSIGNED NOT NULL, turn INT, PRIMARY KEY (live_id, turn), FOREIGN KEY (live_id) REFERENCES Live(id) );
CREATE TABLE Player ( name VARCHAR(50), part ENUM('Vo.', 'Gt.', 'Gt.Vo.', 'Key.', 'Ba.', 'Dr.'), PRIMARY KEY (name, part) );
CREATE TABLE BandMember ( live_id BIGINT UNSIGNED NOT NULL, turn INT, member_name VARCHAR(50), member_part ENUM('Vo.', 'Gt.', 'Gt.Vo.', 'Key.', 'Ba.', 'Dr.'), PRIMARY KEY(live_id, turn, member_name, member_part), FOREIGN KEY (live_id, turn) REFERENCES Band(live_id, turn), FOREIGN KEY (member_name, member_part) REFERENCES Player(name, part) ON UPDATE CASCADE );

# データ挿入
## Live
INSERT INTO Live(name, location, date, performance_fee, equipment_cost) VALUES ('name', 'location', '2022-01-03', 5500, 2000);

## Band
INSERT INTO Band VALUES('name', 1, 1);
INSERT INTO Band VALUES('name2', 1, 2);

## Player
INSERT INTO Player VALUES ('drummer', 'Dr.');
INSERT INTO Player VALUES ('guiterist', 'Gt.');

## BandMember
INSERT INTO BandMember VALUES(1, 1, 'drummer', 'Dr.');
INSERT INTO BandMember VALUES(1, 2, 'guiterist', 'Gt.');

SELECT * FROM Band WHERE live_id IN (SELECT id FROM Live WHERE date = '2022-01-03');
```