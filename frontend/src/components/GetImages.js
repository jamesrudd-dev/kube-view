import React, { useState, useEffect } from 'react';
import axios from 'axios';

axios.defaults.baseURL = 'http://localhost:8080';

const GetImages = ({ cluster, namespace }) => {
  const [images, setImages] = useState([]);
  const [isLoading, setIsLoading] = useState(false);
  const [isVisible, setIsVisible] = useState(false);

  useEffect(() => {
    if (cluster !== undefined && namespace !== undefined) {
      setIsLoading(true);
      setIsVisible(false);
      setImages([])
      axios
      .get(`/deployments/${cluster}/${namespace}`)
      .then((res) => {
        console.log(res);
        setImages(res.data);
      })
      .catch((err) => {
        console.log(err);
      })
      .finally(() => {
          setIsLoading(false);
          setTimeout(() => {
            setIsVisible(true);
          }, 1500);
      })
    }
  }, [cluster, namespace]);

  if (Object.keys(images).length > 0) {
    return (
      <div>  

        {isLoading && <h2 className="text-light">Fetching from {namespace}</h2>}
  
        <div className='container'>
          {images.map((data) => (
  
            <div className='card w-100 text-white bg-dark' key={data.id}>
              <h4 className='card-header'>{data.deploymentName}</h4>
              <p>{data.imageName}</p>
              <p>{data.imageTag}</p>
            </div>
          ))}
        </div>
  
      </div>
    );
  }
  if (Object.keys(images).length === 0) {
    return (
      <div>
        {isLoading && <h2 className="text-light">Fetching Deployments from {namespace}</h2>}
        {isVisible && <h1>No deployments in this namespace.</h1>}
      </div>
    )
  }

};

export default GetImages;