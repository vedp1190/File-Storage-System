import React, { useState } from 'react';

function App() {
  const [selectedFile, setSelectedFile] = useState();
  const [isFilePicked, setIsFilePicked] = useState(false);
  const [files, setFiles] = useState([]);

  const changeHandler = (event) => {
    setSelectedFile(event.target.files[0]);
    setIsFilePicked(true)
  }

  const handleSubmission = () => {
    const formData = new FormData();
    formData.append('file', selectedFile);
    fetch(
      'http://localhost:8080/upload',
      {
        method: 'POST',
        body: formData,
      }
    )
      .then((result => {
        console.log('Success', result);
        setFiles([...files, { name: selectedFile.name, type: selectedFile.type }])
      }))
      .catch((error) => {
        console.error('Error', error);
      })
  };

  return (
    <div className="flex flex-col justify-center items-center h-screen bg-darkbg text-white">
      <input type="file" name="file" onChange={changeHandler}>
      </input>
      <div>
        <button onClick={handleSubmission}>Submit</button>
      </div>
      <ul>
        {
          files.map(file => {
            return <li onClick={() => download(file.name)}>{file.name}</li>
          })
        }
      </ul>
    </div>
  );
}

function download(name) {
  console.log("aaaaa")
  fetch('http://localhost:8080/download?id=' + name)
    .then(response => {
      response.blob().then(blob => {
        let url = window.URL.createObjectURL(blob);
        let a = document.createElement('a');
        a.href = url;
        a.download = name;
        a.click();
      });
    })
    .then((result => {
      console.log('Success', result);
    }))
    .catch((error) => {
      console.error('Error', error);
    })
}

export default App;


/*
import React, { useState } from 'react';

function App() {
  const [selectedFile, setSelectedFile] = useState();
  const [isFilePicked, setIsFilePicked] = useState(false);

  const changeHandler = (event) => {
    setSelectedFile(event.target.files[0]);
    setIsFilePicked(true)
  }

  const handleSubmission = () => {
  };

  return (
    <div className="flex flex-col justify-center items-center h-screen bg-darkbg">
      <div className='text-white'>Your Files</div>
      <div className='w-1/2 h-1/4 bg-white rounded-md'></div>
      <div className='bg-white rounded-lg w-24 h-8 m-5 flex justify-center items-center select-none cursor-pointer'>Upload</div>
      <div className='bg-white'>
        <input type="file" name="file" onChange={changeHandler}>
          {isFilePicked ? (
            <div>
              <p>Filename: {selectedFile.name}</p>
              <p>Filetype: {selectedFile.type}</p>
              <p>Size in bytes: {selectedFile.size}</p>
              <p>
                lastModifiedDate:{' '}
                {selectedFile.lastModifedDate.toLocaleDateString()}
              </p>
            </div>
          ) : (
            <p>Select a file to show details.</p>
          )}
        </input>
        <div>
          <button onClick={handleSubmission}>Submit</button>
        </div>
      </div>
    </div>
  );
}

export default App;
*/
