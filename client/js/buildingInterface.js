var BuildingInterface = React.createClass({
   build: function(type) {
      var i = this.props.selectedI;
      var j = this.props.selectedJ;

      socket.send(JSON.stringify({
         command: "build",
         paramsBuild: {
            type: type,
            planetId: this.props.planet.id,
            i: i,
            j: j
         }
      }));

      this.props.unselectBuilding();
   },

   buildShip: function(type) {
      var i = this.props.selectedI;
      var j = this.props.selectedJ;

      socket.send(JSON.stringify({
         command: "buildShip",
         paramsBuildShip: {
            type: type,
            name: prompt("Ship name", "S.S. Enterprise"),
            planetId: this.props.planet.id,
         }
      }));

      this.props.unselectBuilding();
   },

   render: function () {
      var b = this.props.selectedBuilding;

      if (b == null)
         return null;

      if (b.operational == false)
         return null;

      if (b.type == "") {
         interfaceBody = (
            <div>
               <div>
                  <div className="tile farm" title={buildingTitle("farm")} onClick={this.build.bind(this, "farm")} ></div>
                  Farm
                  <div className="buildingCosts">
                     0 <span className="resource obtanium" title="Obtanium"></span>
                     1k <span className="resource workers" title="Workers"></span>
                     50<span className="resource energy" title="Energy"></span>
                  </div>
               </div>
               <div>
                  <div className="tile generator" title={buildingTitle("generator")} onClick={this.build.bind(this, "generator")} ></div>
                  Generator
                  <div className="buildingCosts">
                     8 <span className="resource obtanium" title="Obtanium"></span>
                     5k <span className="resource workers" title="Workers"></span>
                     0 <span className="resource energy" title="Energy"></span>
                  </div>
               </div>
               <div>
                  <div className="tile vale" title={buildingTitle("vale")} onClick={this.build.bind(this, "vale")} ></div>
                  Vale
                  <div className="buildingCosts">
                     80 <span className="resource obtanium" title="Obtanium"></span>
                     50k <span className="resource workers" title="Workers"></span>
                     2k <span className="resource energy" title="Energy"></span>
                  </div>
               </div>
               <div>
                  <div className="tile nasa" title={buildingTitle("nasa")} onClick={this.build.bind(this, "nasa")} ></div>
                  Nasa
                  <div className="buildingCosts">
                     40 <span className="resource obtanium" title="Obtanium"></span>
                     100k <span className="resource workers" title="Workers"></span>
                     1k <span className="resource energy" title="Energy"></span>
                  </div>
               </div>
            </div>
         );
      } else if (b.type == "nasa") {
         interfaceBody = (
            <div>
               <div>
                  <div className="shipIcon colonizer" onClick={this.buildShip.bind(this,"colonizer")}></div>
                  Colonizer
                  <div className="shipCosts">
                     100 <span className="resource obtanium" title="Obtanium"></span>
                  </div>
               </div>
            </div>
         );
      } else {
         interfaceBody = null;
      }

      return (
         <div>
            <div className="unselector" onClick={this.props.unselectBuilding}></div>

            <div id="buildingInterface">
               <h2>{buildingTitle(b)}</h2>

               {interfaceBody}
            </div>
         </div>
      );
   }
});
