alter table "public"."assignments" drop constraint "assignments_test_id_fkey",
  add constraint "candidates_test_id_fkey"
  foreign key ("test_id")
  references "public"."tests"
  ("id") on update restrict on delete set null;
