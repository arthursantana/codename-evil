var ManagementInterface = React.createClass({
   changePlanetName: function () {
      var self = this;
      var newName = prompt("Planet name", this.props.selectedPlanet.name);

      if (newName)
         socket.send(JSON.stringify({
            command: "changePlanetName",
            paramsChangePlanetName: {
               id: self.props.selectedPlanet.id,
               name: newName
            }
         }));
   },

   render: function () {
      if (this.props.planets == null || this.props.players == null) // no data yet
         return null;

      if (this.props.selectedPlanet == null) {
         return null;
      }

      if (this.props.selectedPlanet.ownerId == -1) {
         interfaceBody = (
            <div id="managementInterface">
               <h2>{this.props.selectedPlanet.name}</h2>
               <div id="playerTag">Nobody owns this shit.</div>
            </div>
         )
      } else {
         interfaceBody = (
            <div id="managementInterface">
               <h2><span>{this.props.selectedPlanet.name}</span><a onClick={this.changePlanetName}><i className="fa fa-edit" aria-hidden={true}></i></a></h2>
               <div id="playerTag"><div id="playerColor"></div>&nbsp;{this.props.players[this.props.selectedPlanet.ownerId].name}</div>

               <ResourceList planet={this.props.selectedPlanet} />
               <WorldMap planet={this.props.selectedPlanet} />

               <ShipList planet={this.props.selectedPlanet} ships={this.props.ships} enterSetVoyageMode={this.props.enterSetVoyageMode} />
            </div>
         )
      }

      return (
         <div>
            <div className="unselector" onClick={this.props.unselectPlanet}></div>

            {interfaceBody}
         </div>
      );
   },

   componentDidUpdate: function () {
      var pC = document.getElementById('playerColor');

      if (pC != null) {
         pC.style.background = this.props.players[this.props.selectedPlanet.ownerId].color;
      }
   }
});
