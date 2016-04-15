var ShipList = React.createClass({
   render: function () {
      var ships = [];

      for (var i = 0; i < this.props.ships.length; i++) {
         var ship = this.props.ships[i]
         if (ship.planetId == this.props.planet.id) {
            newShipIcon = <div key={i} className={"shipIcon "+ship.type} title={ship.name} onClick={this.props.enterSetVoyageMode.bind(null, this.props.ships[i])}></div>;
            ships.push(newShipIcon);
         }
      }

      return (
         <div id="shipList">
            {ships}
         </div>
      );
   }
});
