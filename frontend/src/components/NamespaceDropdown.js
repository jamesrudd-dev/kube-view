import React, { useState, useEffect } from 'react';
import axios from 'axios';
import GetImages from './GetImages';

axios.defaults.baseURL = 'http://localhost:8080';

const NamespaceDropdown = ({ cluster, handleNamespaceChange }) => {
    const [namespace, SetNamespace] = useState([]);
    const [currentNamespace, setCurrentNamespace] = useState("");

    useEffect(() => {
        if (cluster !== undefined) {
            axios
            .get(`/cluster/${cluster}/namespaces`)
            .then((res) => {
                console.log(res);
                SetNamespace(res.data);
                var first = res.data[0]
                console.log(first)
                setCurrentNamespace(first.namespace)
            })
            .catch((err) => {
                console.log(err);
            });
        }
    }, [cluster]);

    if (currentNamespace !== "") {
        return (
            <div>  
                <label className="text-light" htmlFor="namespaceDropdown">Namespace:</label>
                <select className="namespace-drop-down" name="namespaceDropdown" onChange={(namespace) => setCurrentNamespace(namespace.target.value)}>
                    {namespace.map(
                        data => <option key={data.id}>{data.namespace}</option>
                    )}
                </select>

                <br></br>
                <br></br>

                <GetImages cluster={cluster} namespace={currentNamespace}/>
            </div>
        );
    } else {
        return (
            <div>  
                <label className="text-light" htmlFor="namespaceDropdown">Namespace:</label>
                <select className="namespace-drop-down" name="namespaceDropdown" value={currentNamespace} onChange={(namespace) => setCurrentNamespace(namespace.target.value)}>
                    {namespace.map(
                        data => <option key={data.id}>{data.namespace}</option>
                    )}
                </select>
            </div>
        );
    }
};

export default NamespaceDropdown;