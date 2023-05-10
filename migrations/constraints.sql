ALTER TABLE "prushka_beta.UserPrivilege" ADD CONSTRAINT "UserPrivilege_fk0" FOREIGN KEY ("privilege_id") REFERENCES "prushka_beta.Privilege"("id");
ALTER TABLE "prushka_beta.UserPrivilege" ADD CONSTRAINT "UserPrivilege_fk1" FOREIGN KEY ("user_id") REFERENCES "prushka_beta.User"("id");
ALTER TABLE "prushka_beta.UserPrivilege" ADD CONSTRAINT "UserPrivilege_fk2" FOREIGN KEY ("workspace_id") REFERENCES "prushka_beta.Workspace"("id");


ALTER TABLE "prushka_beta.Desk" ADD CONSTRAINT "Desk_fk0" FOREIGN KEY ("workspace_id") REFERENCES "prushka_beta.Workspace"("id");

ALTER TABLE "prushka_beta.Column" ADD CONSTRAINT "Column_fk0" FOREIGN KEY ("desk_id") REFERENCES "prushka_beta.Desk"("id");

ALTER TABLE "prushka_beta.Card" ADD CONSTRAINT "Card_fk0" FOREIGN KEY ("column_id") REFERENCES "prushka_beta.Column"("id");
ALTER TABLE "prushka_beta.Card" ADD CONSTRAINT "Card_fk1" FOREIGN KEY ("assigned") REFERENCES "prushka_beta.User"("id");
ALTER TABLE "prushka_beta.Card" ADD CONSTRAINT "Card_fk2" FOREIGN KEY ("creator") REFERENCES "prushka_beta.User"("id");


ALTER TABLE "prushka_beta.Attachment" ADD CONSTRAINT "Attachment_fk0" FOREIGN KEY ("card_id") REFERENCES "prushka_beta.Card"("id");

ALTER TABLE "prushka_beta.CardsLabel" ADD CONSTRAINT "CardsLabel_fk0" FOREIGN KEY ("card_id") REFERENCES "prushka_beta.Card"("id");
ALTER TABLE "prushka_beta.CardsLabel" ADD CONSTRAINT "CardsLabel_fk1" FOREIGN KEY ("label_id") REFERENCES "prushka_beta.Label"("id");
