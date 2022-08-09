import React from 'react';

class NamespaceDropdown extends React.Component {
    state = {
        values: []
    }
    componentDidMount() {
        let cluster = "epe-kubernetes"
        fetch(`http://localhost:8080/cluster/${cluster}/namespaces`)
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
                     return <option value={obj.id}>{obj.namespace}</option>
                 })
              }</select>
            </div>;
    }
}

export default NamespaceDropdown;