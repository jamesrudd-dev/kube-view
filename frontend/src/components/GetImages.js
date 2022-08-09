import React, { useState, useEffect } from 'react';
import axios from 'axios';

axios.defaults.baseURL = 'http://localhost:8080';

const GetImages = ({ cluster, namespace }) => {
  const [images, setImages] = useState([]);


  useEffect(() => {
    if (cluster !== undefined && namespace !== undefined) {
      const fetchImages = () => {
        axios
        .get(`/deployments/${cluster}/${namespace}`)
        .then((res) => {
          console.log(res);
          setImages(res.data);
        })
        .catch((err) => {
          console.log(err);
        });
      };
      fetchImages();
    }
  }, [cluster, namespace]);


  return (
    <div>  

      <div className='container'>
        {images.map((image) => (

          <div className='card w-100 text-white bg-dark' key={image.id}>
            <h4 className='card-header'>{image.deploymentName}</h4>
            <p>{image.imageName}</p>
            <p>{image.imageTag}</p>
          </div>
        ))}
      </div>

    </div>
  );
};

export default GetImages;