package schema

func init() {
	patches.Register(
		1,
		`
	-- ----------------------------------------------------------------
	-- Name: move_events
	-- Model: model.move_events
	-- ----------------------------------------------------------------
	CREATE TABLE {{ .SchemaName | default "public"}}.move_events (
		height 				bigint NOT NULL,
		parent_state_root	text NOT NULL,
		parent_tip_set		text NOT NULL,
		tip_set 			text
	);
	ALTER TABLE ONLY {{ .SchemaName | default "public"}}.chain_consensus ADD CONSTRAINT chain_consensus_pkey PRIMARY KEY (height, parent_state_root, parent_tip_set);
`)

}
