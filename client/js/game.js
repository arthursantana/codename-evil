var socket;

var lastUpdateTimestamp = 0

var Game = React.createClass({
   getInitialState: function () {
      return {
         players: null,
         planets: null,
         ships: null,
         units: null,
         selectedShip: null,
         selectedPlanet: null
      };
   },

   unselectPlanet: function () {
      this.setState({ selectedPlanet: null });
   },

   setSelectedPlanet: function (planet) {
      this.setState({ selectedPlanet: planet });
   },

   enterSetDestinationMode: function (ship) {
      this.setState({ selectedPlanet: null, selectedShip: ship });
   },
   
   quitSetDestinationMode: function () {
      this.setState({ selectedShip: null });
   },

   render: function () {
      return (
         <div>
            <StarMap planets={this.state.planets} players={this.state.players} ships={this.state.ships} setSelectedPlanet={this.setSelectedPlanet} selectedShip={this.state.selectedShip} quitSetDestinationMode={this.quitSetDestinationMode} />
            <ManagementInterface planets={this.state.planets} players={this.state.players} ships={this.state.ships} units={this.state.units} selectedPlanet={this.state.selectedPlanet} unselectPlanet={this.unselectPlanet} enterSetDestinationMode={this.enterSetDestinationMode} />
         </div>
      );
   },

   componentDidMount: function () {
      var self = this;

      socket = new WebSocket("ws://" + window.location.hostname + ":" + window.location.port + "/ws/");

      socket.onopen = function (event) {
         if (document.cookie == null || document.cookie == 0) {
            var player = {
               name: prompt("Name:", ""),
               color: prompt("Color:", "#005500"),
               points: 0
            }

            socket.send(JSON.stringify(player)); 
         } else {
            socket.send(JSON.stringify({
               name: "___reconnect___",
               color: document.cookie
            })); 
         }

         socket.onmessage = function (event) {
            var answer = JSON.parse(event.data);

            if (answer.type == "dataUpdate") {
               var selectedPlanet = null

               if (answer.timestamp > lastUpdateTimestamp) {
                  lastUpdateTimestamp = answer.timestamp;
               } else {
                  console.log("Error: received data from the past.");
                  return;
               }

               if (self.state.selectedPlanet != null) { // gotta find this planet in the new planet list
                  var selectedPlanetId = self.state.selectedPlanet.id;

                  for (var i = 0; i < answer.planets.length; i++) {
                     if (answer.planets[i].id == selectedPlanetId) {
                        selectedPlanet = answer.planets[i];
                        break;
                     }
                  }
               }

               self.setState({
                  players: answer.players,
                  planets: answer.planets,
                  ships: answer.ships,
                  units: answer.units,
                  selectedPlanet: selectedPlanet
               });
            }
            else { // stores registered userId
               document.cookie = "#" + event.data;
            }
         };
      };
   }
});

ReactDOM.render(React.createElement(Game, {}), document.getElementById("game"));
