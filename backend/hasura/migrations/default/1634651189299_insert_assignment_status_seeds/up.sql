INSERT INTO public.assignment_status (value) VALUES ('sending') ON CONFLICT DO NOTHING;
INSERT INTO public.assignment_status (value) VALUES ('sent') ON CONFLICT DO NOTHING;
INSERT INTO public.assignment_status (value) VALUES ('viewed') ON CONFLICT DO NOTHING;
INSERT INTO public.assignment_status (value) VALUES ('scheduled') ON CONFLICT DO NOTHING;
INSERT INTO public.assignment_status (value) VALUES ('cancelled') ON CONFLICT DO NOTHING;
INSERT INTO public.assignment_status (value) VALUES ('submitted') ON CONFLICT DO NOTHING;
INSERT INTO public.assignment_status (value) VALUES ('inprogress') ON CONFLICT DO NOTHING;
INSERT INTO public.assignment_status (value) VALUES ('missed') ON CONFLICT DO NOTHING;
