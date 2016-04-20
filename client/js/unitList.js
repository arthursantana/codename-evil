var UnitList = React.createClass({
   handleDragStart: function (e) {
      e.dataTransfer.effectAllowed = "move";
      console.log(e.target.dataset.id)
      e.dataTransfer.setData("text/html", e.target.dataset.id);
   },

   render: function () {
      var units = [];

      for (var i = 0; i < this.props.units.length; i++) {
         var unit = this.props.units[i]
         if (unit.planetId == this.props.planet.id) {
            if (unit.ownerId == this.props.planet.ownerId)
               enemyUnitClass = "";
            else
               enemyUnitClass = " enemy";
            newUnitIcon = <div key={i} data-id={unit.id} className={"unitIcon " + unit.type + enemyUnitClass} title={unit.type} draggable="true" onDragStart={this.handleDragStart}></div>;
            units.push(newUnitIcon);
         }
      }

      return (
         <div id="unitList">
            {units}
         </div>
      );
   }
});
