var WorldMap = React.createClass({
   getInitialState: function () {
      return {
         selectedI: -1,
         selectedJ: -1
      };
   },

   unselectBuilding: function () {
      this.setState({
         selectedI: -1,
         selectedJ: -1
      });
   },

   setSelectedBuilding: function (i, j) {
      this.setState({
         selectedI: i,
         selectedJ: j
      });
   },

   render: function () {
      var buildings = this.props.planet.buildings;
      var rows = [];

      for (var i = 0; i < buildings.length; i++) {
         var tiles = [];
         var op;

         for (var j = 0; j < buildings.length; j++) {
            if (buildings[i][j].operational == true)
               op = "";
            else
               op = " notOperational";

            if (buildings[i][j].ticksUntilDone == 0)
               stillBuilding = null;
            else
               stillBuilding = <div className="stillBuilding"></div>;

            tiles.push(<div key={j} className={"tile "+buildings[i][j].type + op} title={buildingTitle(buildings[i][j])} onClick={this.setSelectedBuilding.bind(this, i, j)} >{stillBuilding}</div>);
         }

         rows.push(
            <div className="row" key={i}>
               {tiles}
            </div>
         );
      }

      return (
         <div>
            <BuildingInterface planet={this.props.planet} selectedI={this.state.selectedI} selectedJ={this.state.selectedJ} unselectBuilding={this.unselectBuilding} />

            <div id="worldMap">
               {rows}
            </div>
         </div>
      );
   },
});
