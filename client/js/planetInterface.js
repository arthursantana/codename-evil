var PlanetInterface = React.createClass({
   changePlanetName: function () {
      var self = this;
      var newName = prompt("Planet name", this.props.openPlanet.name);

      if (newName)
         socket.send(JSON.stringify({
            command: "changePlanetName",
            paramsChangePlanetName: {
               id: self.props.openPlanet.id,
               name: newName
            }
         }));
   },

   render: function () {
      if (this.props.players == null) // no data yet
         return null;

      if (this.props.openPlanet == null) {
         return null;
      }

      if (this.props.openPlanet.ownerId == -1) {
         interfaceBody = (
            <div id="planetInterface">
               <h2>{this.props.openPlanet.name}</h2>
               <div id="playerTag">Nobody owns this shit.</div>
            </div>
         )
      } else {
         interfaceBody = (
            <div id="planetInterface">
               <h2><span>{this.props.openPlanet.name}</span><a onClick={this.changePlanetName}><i className="fa fa-edit" aria-hidden={true}></i></a></h2>
               <div id="playerTag"><div id="playerColor"></div>&nbsp;{this.props.players[this.props.openPlanet.ownerId].name}</div>

               <ResourceList planet={this.props.openPlanet} />
               <WorldMap planet={this.props.openPlanet} />

               <UnitList planet={this.props.openPlanet} units={this.props.units} />
               <ShipList planet={this.props.openPlanet} ships={this.props.ships} enterSetDestinationMode={this.props.enterSetDestinationMode} />
            </div>
         )
      }

      return (
         <div>
            <div className="unselector" onClick={this.props.closePlanetInterface}></div>

            {interfaceBody}
         </div>
      );
   },

   componentDidUpdate: function () {
      var pC = document.getElementById('playerColor');

      if (pC != null) {
         pC.style.background = this.props.players[this.props.openPlanet.ownerId].color;
      }
   }
});
