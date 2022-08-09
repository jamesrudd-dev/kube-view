import React, { useState, useEffect } from 'react';
import axios from 'axios';
axios.defaults.baseURL = 'http://localhost:8080';
const GetImages = () => {
  const [images, setImages] = useState([]);

  useEffect(() => {
    fetchImages();
  }, []);

  const fetchImages = () => {
    let cluster = 'epe-kubernetes'
    let namespace = 'dev'
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

  return (
    <div className='container'>
      {images.map((image) => (

        <div className='card w-100 text-white bg-dark' key={image.id}>
          <h4 className='card-header'>{image.deploymentName}</h4>
          <p>{image.imageName}</p>
          <p>{image.imageTag}</p>
        </div>

      ))}
    </div>
  );
};

export default GetImages;