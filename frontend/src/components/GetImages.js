import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { TailSpin } from  'react-loader-spinner';

axios.defaults.baseURL = 'http://localhost:8080';

const GetImages = ({ cluster, namespace, refresh }) => {
  const [images, setImages] = useState([]);
  const [isLoading, setIsLoading] = useState(false);
  const [isVisible, setIsVisible] = useState(false);

  useEffect(() => {
    if (cluster !== undefined && namespace !== "" && !refresh) {
      setImages([])
      setIsLoading(true);
      setIsVisible(false);
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
        }, 500);
      })
    }
  }, [cluster, namespace, refresh]);

  if (Object.keys(images).length > 0) {
    return (
      <div>  

        <div className='loading-spinner'>
          {isLoading && <TailSpin
              color = 'white'
              ariaLabel = 'tailspin-loading'     
          />}
        </div>
  
        <div className='grid-container'>
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

        <div className='loading-spinner'>
          {isLoading && <TailSpin
              color = 'white'
              ariaLabel = 'tailspin-loading'     
          />}
        </div>

        <div>
          {isVisible && <h2>No deployments in this namespace</h2>}
        </div>

      </div>
    )
  }

};

export default GetImages;