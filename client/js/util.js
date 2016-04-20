function buildingTitle (building) {
   switch (building.type) {
      case "": return "Click to Build Something";
      case "hq": return "Planet Headquarters";
      case "farm": return "Old McDonald's Farm";
      case "generator": return "Monty Burns Nuclear Energy";
      case "nasa": return "Not Actually a Space Agency";
      case "vale": return "Vale Obtanium Inc.";
      case "lockheed": return "Lockheed Martian Inc.";
   }
}

function commalizer(x) {
   return x.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ",");
}

function capitalize(string) {
   return string.charAt(0).toUpperCase() + string.slice(1);
}
