import React from 'react';

class ClusterDropdown extends React.Component {
    state = {
        values: []
    }
    componentDidMount() {
        fetch(`http://localhost:8080/cluster/list`)
        .then(function(res) {
            return res.json();
        }).then((json)=> {
            this.setState({
               values: json
            })
        });
    }
    render(){
        return <div className="drop-down">
              <select>{
                 this.state.values.map((obj) => {
                     return <option value={obj.id}>{obj.cluster}</option>
                 })
              }</select>
            </div>;
    }
}

export default ClusterDropdown;