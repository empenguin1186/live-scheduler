# データベース設計

## Live テーブル
- id: ライブのID(Auto Increment, 主キー)
- name: ライブ名
- location: ライブの開催場所
- date: ライブ開催日
- performanceFee: 一人当たりの出演費
- equipmentCost: １バンドあたりの機材費

## Band テーブル
- name: バンド名
- live_id: ライブID(主キー) 外部キーとして Live テーブルの id カラムを参照する
- order: 出演順(主キー)

## BandMember テーブル
- live_id: ライブID(主キー) Live テーブルの id カラムを外部キー
- order: 出演順(主キー) Band テーブルの order カラムを外部キー
- member_name: メンバーの名前(主キー) Member テーブルの name カラムを外部キー
- member_part: 担当パート(主キー) Member テーブルの part カラムを外部キー

## Player テーブル
- name: メンバーの名前(主キー)
- part: 担当パート(主キー) enum 型