import React, { useState, useEffect } from 'react';
import axios from 'axios';

axios.defaults.baseURL = 'http://localhost:8080';

const NamespaceDropdown = ({ cluster, handleNamespaceChange }) => {
    const [namespace, SetNamespace] = useState([]);

    useEffect(() => {
        if (cluster !== undefined) {
            axios
            .get(`/cluster/${cluster}/namespaces`)
            .then((res) => {
                console.log(res);
                SetNamespace(res.data);
                handleNamespaceChange(res.data[0]["namespace"])
            })
            .catch((err) => {
                console.log(err);
            });
        }
    }, [cluster, handleNamespaceChange]);

    return (
        <div>  
            <label className="text-light" htmlFor="namespaceDropdown">Namespace:</label>
            <select className="namespace-drop-down" name="namespaceDropdown" onChange={(namespace) => handleNamespaceChange(namespace.target.value)}>
                {namespace.map(
                    data => <option key={data.id}>{data.namespace}</option>
                )}
            </select>
        </div>
    );
};

export default NamespaceDropdown;