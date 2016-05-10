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
         openPlanets: null
      };
   },

   closePlanetInterface: function () {
      this.setState({ openPlanets: null });
   },

   openPlanetInterface: function (planets) {
      this.setState({ openPlanets: planets });
   },

   enterSetDestinationMode: function (ship) {
      this.setState({ openPlanets: null, selectedShip: ship });
   },
   
   quitSetDestinationMode: function () {
      this.setState({ selectedShip: null });
   },

   render: function () {
      if (this.state.openPlanets == null)
         openInterface = null;
      else if (this.state.openPlanets.size == 1) {
         var openPlanetIds = [v for (v of this.state.openPlanets.values())];
         openInterface = <PlanetInterface openPlanet={this.state.planets[openPlanetIds[0]]} players={this.state.players} ships={this.state.ships} units={this.state.units} closePlanetInterface={this.closePlanetInterface} enterSetDestinationMode={this.enterSetDestinationMode} />;
      } else {
         openInterface = <MultiPlanetInterface openPlanets={this.state.openPlanets} planets={this.state.planets} players={this.state.players} ships={this.state.ships} units={this.state.units} closePlanetInterface={this.closePlanetInterface} enterSetDestinationMode={this.enterSetDestinationMode} />;
      }
      return (
         <div>
            <StarMap planets={this.state.planets} players={this.state.players} ships={this.state.ships} openPlanetInterface={this.openPlanetInterface} selectedShip={this.state.selectedShip} quitSetDestinationMode={this.quitSetDestinationMode} />
            {openInterface}
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
               if (answer.timestamp > lastUpdateTimestamp) {
                  lastUpdateTimestamp = answer.timestamp;
               } else {
                  console.log("Error: received data from the past.");
                  return;
               }

               self.setState({
                  players: answer.players,
                  planets: answer.planets,
                  ships: answer.ships,
                  units: answer.units
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
