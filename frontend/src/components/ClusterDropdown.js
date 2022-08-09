import GetImages from './GetImages';
import NamespaceDropdown from './NamespaceDropdown';
import ClusterRefresh from './ClusterRefresh';
import React, { useState, useEffect, useCallback } from 'react';
import axios from 'axios';

axios.defaults.baseURL = 'http://localhost:8080';

const ClusterDropdown = () => {
    const [clusters, SetCluster] = useState([]);
    const [currentCluster, setCurrentCluster] = useState();
    const [currentNamespace, setCurrentNamespace] = useState();

    const handleClusterChange = cluster => {
        setCurrentCluster(cluster)
    };

    const handleNamespaceChange = useCallback(
        namespace => {
            setCurrentNamespace(namespace)
        }, 
    []);

    useEffect(() => {
        axios
        .get(`/cluster/list`)
        .then((res) => {
            console.log(res);
            SetCluster(res.data);
        })
        .catch((err) => {
            console.log(err);
        });
    }, []);

    return (
        <div>
            <div className="row">

                <div className="column">
                    <label className="text-light" htmlFor="clusterDropdown">Cluster:</label>
                    <select className="cluster-drop-down" name="clusterDropdown" onChange={(cluster) => handleClusterChange(cluster.target.value)}>
                        <option value="" disabled>Select Cluster</option>
                        {clusters.map(
                            data => <option key={data.id}>{data.cluster}</option>
                        )}
                    </select>
                </div>

                <div className="column">
                <ClusterRefresh />
                </div>

                <div className="column">
                <NamespaceDropdown cluster={currentCluster} handleNamespaceChange={handleNamespaceChange}/>
                </div>
            </div>

            <br></br>

            {/* <GetImages cluster={currentCluster} namespace={currentNamespace}/> */}

        </div>
    );
};

export default ClusterDropdown;