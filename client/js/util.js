function buildingTitle (building) {
   switch (building.type) {
      case "": return "Click to Build Something";
      case "hq": return "Planet Headquarters";
      case "farm": return "Old McDonald's Farm";
      case "generator": return "Monty Burns Nuclear Energy";
      case "nasa": return "Not Actually a Space Agency";
      case "vale": return "Vale Obtanium Inc.";
   }
}
