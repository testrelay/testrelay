alter table "public"."assignments" drop constraint "candidates_test_id_fkey",
  add constraint "assignments_test_id_fkey"
  foreign key ("test_id")
  references "public"."tests"
  ("id") on update cascade on delete cascade;
