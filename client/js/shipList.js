var ShipList = React.createClass({
   handleDragOver: function (e) {
      if (e.preventDefault) {
         e.preventDefault();
      }

      e.dataTransfer.dropEffect = "move";

      return false;
   },

   handleDrop: function (e) {
      if (e.stopPropagation) {
         e.stopPropagation(); // stops the browser from redirecting.
      }

      var unitId = Number(e.dataTransfer.getData("text/html"));
      var shipId = Number(e.target.dataset.id);

      if (this.props.ships[shipId].type == "trojan") {
         socket.send(JSON.stringify({
            command: "boardShip",
            paramsBoardShip: {
               unitId: unitId,
               shipId: shipId
            }
         }));

         return false;
      } else {
         return false;
      }
   },

   render: function () {
      var ships = [];

      for (var i = 0; i < this.props.ships.length; i++) {
         var ship = this.props.ships[i]
         if (ship.planetId == this.props.planet.id) {
            if (ship.ownerId == this.props.planet.ownerId)
               enemyShipClass = "";
            else
               enemyShipClass = " enemy";
            newShipIcon = <div key={i} data-id={ship.id} className={"shipIcon " + ship.type + enemyShipClass} title={ship.name + " (" + ship.unitSpace + " slots available)"} onClick={this.props.enterSetVoyageMode.bind(null, this.props.ships[i])} onDragOver={this.handleDragOver} onDrop={this.handleDrop}></div>;
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
