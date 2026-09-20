export const PRESENTATION_LAYERS = [
  "presentation.world_map",
  "presentation.portfolio_map",
  "presentation.site_view",
  "presentation.operational_detail",
] as const;

export type PresentationLayer = (typeof PRESENTATION_LAYERS)[number];

export type FoundationArea = {
  id: string;
  owner: string;
  state: "foundation";
};

export const FOUNDATION_AREAS: FoundationArea[] = [
  { id: "identity.account", owner: "services/api", state: "foundation" },
  { id: "accounting.ledger", owner: "services/api", state: "foundation" },
  { id: "property.real_estate", owner: "packages/content", state: "foundation" },
  { id: "moderation.policy", owner: "services/realtime", state: "foundation" },
];
