import { imageState } from '../hooks/useImageProcessor';
import styles from './ImageForm.module.css';
import { fetchWithAuth } from '../utils/api';

export function ImageForm() {
  const handleSubmit = async (event: Event) => {
    event.preventDefault();
    
    imageState.value = {
      ...imageState.value,
      errorMessage: null,
      progress: 0,
      downloadUrl: null
    };

    if (!imageState.value.imageFile) {
      imageState.value = {
        ...imageState.value,
        errorMessage: 'Please select an image file.'
      };
      return;
    }

    const formData = new FormData();
    formData.append('image', imageState.value.imageFile);
    formData.append('width', imageState.value.width);
    formData.append('height', imageState.value.height);
    formData.append('maintainAspectRatio', imageState.value.maintainAspectRatio.toString());
    formData.append('quality', imageState.value.quality.toString());
    if (imageState.value.maxSize) formData.append('maxSize', imageState.value.maxSize);

    try {
      const response = await fetchWithAuth('/api/resize', {
        method: 'POST',
        body: formData,
      });

      if (!response.ok) {
        throw new Error('Failed to submit the image. Please try again.');
      }

      const data = await response.json();
      imageState.value = { ...imageState.value, jobId: data.job_id };
      pollProgress();
    } catch (error) {
      if (error instanceof Error) {
        imageState.value = { ...imageState.value, errorMessage: error.message };
      }
    }
  };

  const pollProgress = () => {
    const interval = setInterval(async () => {
      try {
        const response = await fetchWithAuth(`/api/progress/${imageState.value.jobId}`);
        const data = await response.json();
        
        imageState.value = { ...imageState.value, progress: data.progress };

        if (data.progress === 100) {
          clearInterval(interval);
          imageState.value = { ...imageState.value, downloadUrl: data.download_url };
          updateProcessedImageData(data.download_url);
        }
      } catch (error) {
        clearInterval(interval);
        imageState.value = {
          ...imageState.value,
          errorMessage: 'Failed to fetch progress. Please try again.'
        };
      }
    }, 2000);
  };

  const updateProcessedImageData = async (downloadUrl: string) => {
    const img = new Image();
    img.onload = () => {
      imageState.value = {
        ...imageState.value,
        processedWidth: img.width,
        processedHeight: img.height
      };
    };
    img.src = downloadUrl;

    const response = await fetch(downloadUrl);
    const blob = await response.blob();
    const processedSize = (blob.size / 1024).toFixed(2) + ' KB';
    imageState.value = { ...imageState.value, processedSize };
  };

  return (
    <form onSubmit={handleSubmit} class={styles.form}>
      <div>
        <label htmlFor="width">Width:</label>
        <input
          type="number"
          id="width"
          value={imageState.value.width}
          onInput={(e) => imageState.value = { ...imageState.value, width: e.currentTarget.value }}
        />
      </div>

      <div>
        <label htmlFor="height">Height:</label>
        <input
          type="number"
          id="height"
          value={imageState.value.height}
          onInput={(e) => imageState.value = { ...imageState.value, height: e.currentTarget.value }}
        />
      </div>

      <div>
        <label>
          <input
            type="checkbox"
            checked={imageState.value.maintainAspectRatio}
            onChange={(e) => imageState.value = { ...imageState.value, maintainAspectRatio: e.currentTarget.checked }}
          />
          Maintain Aspect Ratio
        </label>
      </div>

      <div>
        <label htmlFor="quality">Quality (1-100):</label>
        <input
          type="number"
          id="quality"
          min="1"
          max="100"
          value={imageState.value.quality}
          onInput={(e) => imageState.value = { ...imageState.value, quality: Number(e.currentTarget.value) }}
        />
      </div>

      <div>
        <label htmlFor="maxSize">Max Size (KB):</label>
        <input
          type="number"
          id="maxSize"
          value={imageState.value.maxSize}
          onInput={(e) => imageState.value = { ...imageState.value, maxSize: e.currentTarget.value }}
        />
      </div>

      <button type="submit">Resize</button>
    </form>
  );
}