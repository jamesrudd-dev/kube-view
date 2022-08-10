import React, { useState, useEffect } from 'react';
import axios from 'axios';
import GetImages from './GetImages';

axios.defaults.baseURL = 'http://localhost:8080';

const ClusterSelectionMenu = () => {
    const [clusters, SetCluster] = useState([]);
    const [currentCluster, setCurrentCluster] = useState();
    const [namespace, SetNamespace] = useState([]);
    const [currentNamespace, setCurrentNamespace] = useState("");
    const [postData, setPostData] = useState();
    const [isLoading, setIsLoading] = useState(false);

    // Used for cluster list
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
  
    // Used for Cluster Refresh Button
    const clusterRefresh = async () => {
      if (currentCluster !== undefined) {
        setIsLoading(true);
        try {
            const {postData} = await axios.post(
            `http://localhost:8080/cluster/${currentCluster}/refresh`,
            {
                headers: {
                'Content-Type': 'application/json',
                Accept: 'application/json',
                },
            },
            );

            console.log(JSON.stringify(postData, null, 4));

            setPostData(postData);
        } catch (err) {
            console.log(err)
        } finally {
            setIsLoading(false);
        }
      }
    };
    console.log(postData);

    // Used for updating namespaces
    useEffect(() => {
        SetNamespace([])
        if (currentCluster !== undefined) {
            axios
            .get(`/cluster/${currentCluster}/namespaces`)
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
    }, [currentCluster]);


    return (
        <div>
            <div className="row mx-0">
                <div className="column form-floating">

                    <select className="form-select form-select-lg mb-3" id="floatingClusterSelect" aria-label="Cluster Selection" defaultValue="" onChange={(cluster) => setCurrentCluster(cluster.target.value)}>
                        <option value="" disabled></option>
                        {clusters.map(
                            data => <option key={data.id}>{data.cluster}</option>
                        )}
                    </select>
                    <label id="form-select-label" htmlFor="floatingClusterSelect">Cluster Selection</label>

                </div>

                <div className="column">
                    <div>

                        <button className="btn btn-outline-light btn-lg m-2" onClick={!isLoading ? clusterRefresh : null }>
                            {isLoading ? 'Refreshing...' : 'Refresh'}
                        </button>

                    </div>
                </div>

                <div className="column form-floating">

                    <select className="form-select form-select-lg mb-3" id="floatingNamespaceSelect" aria-label="Namespace Selection" onChange={(namespace) => setCurrentNamespace(namespace.target.value)}>
                        {namespace.map(
                                data => <option key={data.id}>{data.namespace}</option>
                        )}
                    </select>
                    <label id="form-select-label" htmlFor="floatingClusterSelect">Namespace Selection</label>
                    
                </div>

            </div>

            <GetImages cluster={currentCluster} namespace={currentNamespace} refreshing={isLoading}/>

        </div>
    );
};

export default ClusterSelectionMenu;