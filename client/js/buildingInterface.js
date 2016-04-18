var BuildingInterface = React.createClass({
   build: function(type) {
      socket.send(JSON.stringify({
         command: "build",
         paramsBuild: {
            type: type,
            planetId: this.props.planet.id,
            i: this.props.selectedI,
            j: this.props.selectedJ
         }
      }));

      this.props.unselectBuilding();
   },

   buildShip: function(type) {
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

   sellBuilding: function() {
      socket.send(JSON.stringify({
         command: "sellBuilding",
         paramsSellBuilding: {
            planetId: this.props.planet.id,
            i: this.props.selectedI,
            j: this.props.selectedJ
         }
      }));

      this.props.unselectBuilding();
   },

   render: function () {
      var b = this.props.selectedBuilding;

      if (b == null)
         return null;

      if (b.operational == false) {
         interfaceBody = null
      }
      else {
         if (b.type == "") {
            interfaceBody = (
               <div>
                  <div>
                     <div className="tile farm" title={buildingTitle("farm")} onClick={this.build.bind(this, "farm")} ></div>
                     Farm
                     <div className="buildingCosts">
                        1k <span className="resource workers" title="Workers"></span>
                        100<span className="resource energy" title="Energy"></span>
                        0 <span className="resource obtanium" title="Obtanium"></span>
                     </div>
                  </div>
                  <div>
                     <div className="tile generator" title={buildingTitle("generator")} onClick={this.build.bind(this, "generator")} ></div>
                     Generator
                     <div className="buildingCosts">
                        10k <span className="resource workers" title="Workers"></span>
                        0 <span className="resource energy" title="Energy"></span>
                        200 <span className="resource obtanium" title="Obtanium"></span>
                     </div>
                  </div>
                  <div>
                     <div className="tile vale" title={buildingTitle("vale")} onClick={this.build.bind(this, "vale")} ></div>
                     Vale
                     <div className="buildingCosts">
                        50k <span className="resource workers" title="Workers"></span>
                        2k <span className="resource energy" title="Energy"></span>
                        800 <span className="resource obtanium" title="Obtanium"></span>
                     </div>
                  </div>
                  <div>
                     <div className="tile nasa" title={buildingTitle("nasa")} onClick={this.build.bind(this, "nasa")} ></div>
                     Nasa
                     <div className="buildingCosts">
                        100k <span className="resource workers" title="Workers"></span>
                        1k <span className="resource energy" title="Energy"></span>
                        400 <span className="resource obtanium" title="Obtanium"></span>
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
                        1k <span className="resource workers" title="Workers"></span>
                        1k <span className="resource cattle" title="Cattle"></span>
                        1k <span className="resource obtanium" title="Obtanium"></span>
                     </div>
                  </div>
               </div>
            );
         } else {
            interfaceBody = null;
         }
      }

      if (b.type != "" && b.type != "hq") {
         sellButton = <a onClick={this.sellBuilding}><i className="fa fa-trash-o" aria-hidden={true}></i></a>;
      } else {
         sellButton = null;
      }

      return (
         <div>
            <div className="unselector" onClick={this.props.unselectBuilding}></div>

            <div id="buildingInterface">
               <h2><span>{buildingTitle(b)}</span>{sellButton}</h2>

               {interfaceBody}
            </div>
         </div>
      );
   }
});
