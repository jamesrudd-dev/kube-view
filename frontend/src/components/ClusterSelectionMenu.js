import NamespaceDropdown from './NamespaceDropdown';
import React, { useState, useEffect } from 'react';
import axios from 'axios';

axios.defaults.baseURL = 'http://localhost:8080';

const ClusterSelectionMenu = () => {
    const [clusters, SetCluster] = useState([]);
    const [currentCluster, setCurrentCluster] = useState();
    const [postData, setPostData] = useState();
    const [isLoading, setIsLoading] = useState(false);
    const [err, setErr] = useState('');

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
          setErr(err.message);
        } finally {
          setIsLoading(false);
          
        }
      }
    };
    console.log(postData);


    return (
        <div>
            <div className="row">
                <div className="column">
                    <label className="text-light" htmlFor="clusterDropdown">Cluster:</label>
                    <select className="cluster-drop-down" name="clusterDropdown" defaultValue="" onChange={(cluster) => setCurrentCluster(cluster.target.value)}>
                        <option value="" disabled>Select Cluster</option>
                        {clusters.map(
                            data => <option key={data.id}>{data.cluster}</option>
                        )}
                    </select>
                </div>

                <div className="column">
                    <div>
                        {err && <h2>{err}</h2>}

                        <button className="btn btn-outline-light" onClick={clusterRefresh}>Refresh</button>

                        {isLoading && <h2>Loading...</h2>}

                    </div>
                </div>

            </div>

            <br></br>

            <NamespaceDropdown cluster={currentCluster}/>

        </div>
    );
};

export default ClusterSelectionMenu;