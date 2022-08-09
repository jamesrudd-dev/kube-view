import React, { useState } from 'react';
import axios from 'axios';

const ClusterRefresh = ({ cluster }) => {
  const [data, setData] = useState();
  const [isLoading, setIsLoading] = useState(false);
  const [err, setErr] = useState('');

  const handleClick = async () => {
    if (cluster !== undefined) {
      setIsLoading(true);
      try {
        const {data} = await axios.post(
          `http://localhost:8080/cluster/${cluster}/refresh`,
          {
            headers: {
              'Content-Type': 'application/json',
              Accept: 'application/json',
            },
          },
        );
  
        console.log(JSON.stringify(data, null, 4));
  
        setData(data);
      } catch (err) {
        setErr(err.message);
      } finally {
        setIsLoading(false);
      }
    }

  };

  console.log(data);

  return (
    <div>
      {err && <h2>{err}</h2>}

      <button className="btn btn-outline-light" onClick={handleClick}>Refresh</button>

      {isLoading && <h2>Loading...</h2>}

    </div>
    
  );
};

export default ClusterRefresh;