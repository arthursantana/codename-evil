var MultiPlanetInterface = React.createClass({
   render: function () {
      var openPlanetIds = [...this.props.openPlanets];

      return (
         <div>
            <div className="unselector" onClick={this.props.closePlanetInterface}></div>

            <div id="multiPlanetInterface">
               <h2><span>Many planets</span></h2>

               <div id="playerTag"><div id="playerColor"></div>&nbsp;{this.props.players[this.props.planets[openPlanetIds[0]].ownerId].name}</div>
               BLARGH
            </div>
         </div>
      );
   },

   componentDidUpdate: function () {
      var pC = document.getElementById('playerColor');
      var openPlanetIds = [...this.props.openPlanets.values()];

      if (pC != null) {
         console.log(this.props.players[this.props.planets[openPlanetIds[0]].ownerId].color);
         pC.style.background = this.props.players[this.props.planets[openPlanetIds[0]].ownerId].color;
      }
   }
});
