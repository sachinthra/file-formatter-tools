<script>
  import { onMount } from 'svelte';

  let imageFile = null;
  let width = '';
  let height = '';
  let maintainAspectRatio = true;
  let quality = 80;
  let maxSize = '';
  let jobId = null;
  let progress = 0;
  let downloadUrl = null;
  let errorMessage = null;
  let isDragging = false;

  // Original image data
  let originalWidth = null;
  let originalHeight = null;
  let originalSize = null;

  // Processed image data
  let processedWidth = null;
  let processedHeight = null;
  let processedSize = null;

  function handleFileChange(event) {
    const file = event.target.files[0];
    processFile(file);
  }

  function handleDragOver(event) {
    event.preventDefault();
    isDragging = true;
  }

  function handleDragLeave() {
    isDragging = false;
  }

  function handleDrop(event) {
    event.preventDefault();
    isDragging = false;

    const file = event.dataTransfer.files[0];
    if (file) {
      processFile(file);
    }
  }

  function processFile(file) {
    if (file) {
      imageFile = file;
      originalSize = (file.size / 1024).toFixed(2) + ' KB'; // File size in KB
      maxSize = originalSize; 

      // Load the image to get its dimensions
      const img = new Image();
      img.onload = () => {
        originalWidth = img.width;
        originalHeight = img.height;

        // Set default width and height in the form
        width = img.width.toString();
        height = img.height.toString();

      };
      img.src = URL.createObjectURL(file);
    }
  }

  async function handleSubmit() {
    errorMessage = null;
    progress = 0;
    downloadUrl = null;

    if (!imageFile) {
      errorMessage = 'Please select an image file.';
      return;
    }

    const formData = new FormData();
    formData.append('image', imageFile);
    formData.append('width', width);
    formData.append('height', height);
    formData.append('maintainAspectRatio', maintainAspectRatio.toString());
    formData.append('quality', quality.toString());
    if (maxSize) formData.append('maxSize', maxSize.toString());

    try {
      const response = await fetch('/api/resize', {
        method: 'POST',
        body: formData,
      });

      if (!response.ok) {
        throw new Error('Failed to submit the image. Please try again.');
      }

      const data = await response.json();
      jobId = data.job_id;

      // Poll for progress
      pollProgress();
    } catch (error) {
      errorMessage = error.message;
    }
  }

  async function pollProgress() {
    const interval = setInterval(async () => {
      try {
        const response = await fetch(`/api/progress/${jobId}`);
        const data = await response.json();
        progress = data.progress;

        if (progress === 100) {
          clearInterval(interval);
          downloadUrl = data.download_url;

          // Fetch processed image dimensions and size
          const img = new Image();
          img.onload = () => {
            processedWidth = img.width;
            processedHeight = img.height;
          };
          img.src = downloadUrl;

          // Fetch processed image size
          const response = await fetch(downloadUrl);
          const blob = await response.blob();
          processedSize = (blob.size / 1024).toFixed(2) + ' KB';
        }
      } catch (error) {
        clearInterval(interval);
        errorMessage = 'Failed to fetch progress. Please try again.';
      }
    }, 2000);
  }
</script>

<main>
  <h1>Single Image Resize</h1>
  {#if errorMessage}
    <p class="error">{errorMessage}</p>
  {/if}
  <div
    class="drop-zone"
    role="button"
    tabindex="0"
    on:dragover|preventDefault={handleDragOver}
    on:dragleave={handleDragLeave}
    on:drop={handleDrop}
    class:is-dragging={isDragging}
  >
    {#if imageFile}
      <p>Selected File: {imageFile.name}</p>
    {:else}
      <p>Drag and drop an image here, or click to select one.</p>
    {/if}
    <input type="file" id="image" accept="image/*" on:change={handleFileChange} style="display: none;" />
  </div>

  {#if originalWidth && originalHeight && originalSize}
    <div class="original-data">
      <h3>Original Image Data</h3>
      <p>Width: {originalWidth}px</p>
      <p>Height: {originalHeight}px</p>
      <p>Size: {originalSize}</p>
    </div>
  {/if}

  <form on:submit|preventDefault={handleSubmit}>
    <div>
      <label for="width">Width:</label>
      <input type="number" id="width" bind:value={width} />
    </div>

    <div>
      <label for="height">Height:</label>
      <input type="number" id="height" bind:value={height} />
    </div>

    <div>
      <label>
        <input type="checkbox" bind:checked={maintainAspectRatio} />
        Maintain Aspect Ratio
      </label>
    </div>

    <div>
      <label for="quality">Quality (1-100):</label>
      <input type="number" id="quality" min="1" max="100" bind:value={quality} />
    </div>

    <div>
      <label for="maxSize">Max Size (KB):</label>
      <input type="number" id="maxSize" bind:value={maxSize} />
    </div>

    <button type="submit">Resize</button>
  </form>

  {#if progress > 0}
    <p>Progress: {progress}%</p>
  {/if}

  {#if downloadUrl}
    <div class="processed-data">
      <h3>Processed Image Data</h3>
      {#if processedWidth && processedHeight && processedSize}
        <p>Width: {processedWidth}px</p>
        <p>Height: {processedHeight}px</p>
        <p>Size: {processedSize}</p>
      {/if}
      <p>
        <a href={downloadUrl} target="_blank" rel="noopener noreferrer">Download Resized Image</a>
      </p>
    </div>
  {/if}
</main>

<style>
  main {
    font-family: Arial, sans-serif;
    max-width: 600px;
    margin: 0 auto;
    padding: 1rem;
  }

  .drop-zone {
    border: 2px dashed #007bff;
    padding: 1rem;
    text-align: center;
    margin-bottom: 1rem;
    cursor: pointer;
  }

  .drop-zone.is-dragging {
    background-color: #e9f5ff;
  }

  .original-data,
  .processed-data {
    margin-top: 1rem;
    padding: 1rem;
    border: 1px solid #ddd;
    background-color: #242424;
  }

  form div {
    margin-bottom: 1rem;
  }

  label {
    display: block;
    margin-bottom: 0.5rem;
  }

  input[type="file"],
  input[type="number"],
  button {
    width: 100%;
    padding: 0.5rem;
    font-size: 1rem;
  }

  button {
    background-color: #007bff;
    color: white;
    border: none;
    cursor: pointer;
  }

  button:hover {
    background-color: #0056b3;
  }
</style>