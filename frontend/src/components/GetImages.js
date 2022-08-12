import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { TailSpin } from  'react-loader-spinner';

const GetImages = ({ cluster, namespace, refreshing }) => {
  const [images, setImages] = useState([]);
  const [isLoading, setIsLoading] = useState(false);
  const [isVisible, setIsVisible] = useState(false);
  const [filteredResults, setFilteredResults] = useState([]);
  const [searchInput, setSearchInput] = useState('');

  useEffect(() => {
    if (cluster !== undefined && namespace !== "" && !refreshing) {
      setSearchInput('');
      setFilteredResults([]);
      setImages([]);
      setIsLoading(true);
      setIsVisible(false);
      axios
      .get(`${process.env.PUBLIC_URL}/deployments/${cluster}/${namespace}`)
      .then((res) => {
        console.log(res);
        setImages(res.data);
      })
      .catch((err) => {
        console.log(err);
      })
      .finally(() => {
        const timer =setTimeout(() => {
          setIsLoading(false);
          setIsVisible(true);
        }, 500);
        return () => clearTimeout(timer);
      })
    }
  }, [cluster, namespace, refreshing]);

  const searchItems = (searchValue) => {
    setSearchInput(searchValue)
    if (searchInput !== '') {
        const filteredData = images.filter((item) => {
            return Object.values(item["deploymentName"]).join('').toLowerCase().includes(searchInput.toLowerCase())
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

        <div className='row'>
          <div className='col-md-8 offset-md-4'>
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

        <div className="d-flex justify-content-center mt-5">
          <div className="p-2"></div>
          <div className="p-2">          
            {isLoading && <TailSpin
                color = 'white'
                ariaLabel = 'tailspin-loading'     
            />}
          </div>
          <div className="p-2"></div>
        </div>

        <div>
          {isVisible && <h2>No deployments in this namespace</h2>}
        </div>

      </div>
    )
  }

};

export default GetImages;