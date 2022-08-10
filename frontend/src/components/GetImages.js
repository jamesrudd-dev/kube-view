import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { TailSpin } from  'react-loader-spinner';

axios.defaults.baseURL = 'http://localhost:8080';

const GetImages = ({ cluster, namespace, refresh }) => {
  const [images, setImages] = useState([]);
  const [isLoading, setIsLoading] = useState(false);
  const [isVisible, setIsVisible] = useState(false);
  const [filteredResults, setFilteredResults] = useState([]);
  const [searchInput, setSearchInput] = useState('');

  useEffect(() => {
    if (cluster !== undefined && namespace !== "" && !refresh) {
      setSearchInput('');
      setFilteredResults([]);
      setImages([]);
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

  const searchItems = (searchValue) => {
    setSearchInput(searchValue)
    if (searchInput !== '') {
        const filteredData = images.filter((item) => {
            return Object.values(item).join('').toLowerCase().includes(searchInput.toLowerCase())
        })
        setFilteredResults(filteredData)
    }
    else{
        setFilteredResults(images)
    }
}

  if (Object.keys(images).length > 0) {
    return (
      <div>  

        <div className='loading-spinner'>
          {isLoading && <TailSpin
              color = 'white'
              ariaLabel = 'tailspin-loading'     
          />}
        </div>

        <div class='row'>
          <div class='col-md-8 offset-md-4'>
            <input
              type="search" id="image-search-bar" className="form-control mt-3"
              icon='search' placeholder='Search Images...'
              onChange={(e) => searchItems(e.target.value)} />
            <label className="form-label" htmlFor="image-search-bar">Search</label>
          </div>
        </div>

        <div className='grid-container'>
            {searchInput.length > 1 ? (
              filteredResults.map((data) => {
                return (
                  <div className='card w-100 text-white bg-dark' key={data.id}>
                    <h4 className='card-header'>{data.deploymentName}</h4>
                    <p>{data.imageName}</p>
                    <p>{data.imageTag}</p>
                  </div>
                )
              })
            ) : (
              images.map((data) => {
                return (
                  <div className='card w-100 text-white bg-dark' key={data.id}>
                    <h4 className='card-header'>{data.deploymentName}</h4>
                    <p>{data.imageName}</p>
                    <p>{data.imageTag}</p>
                  </div>
                )
              })
            )}
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