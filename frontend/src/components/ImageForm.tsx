import { imageState } from '../hooks/useImageProcessor';
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

      // const transformedUrl = data.download_url ? transformMinioUrl(data.download_url) : null;
      // console.log('[DEBUG] Original URL:', data.download_url);
      // console.log('[DEBUG] Transformed URL:', transformedUrl);

      imageState.value = {
        ...imageState.value,
        jobId: data.job_id,
        downloadUrl: data.download_url,
        objectName: data.object_name
      };

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
        console.log('[DEBUG] Progress data:', data);

        imageState.value = { ...imageState.value, progress: data.progress };

        if (data.progress === 100) {
          clearInterval(interval);
          if (imageState.value.downloadUrl) {
            updateProcessedImageData(imageState.value.downloadUrl);
          }
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
    try {
      // Use fetchWithAuth for both image and blob requests
      const imgResponse = await fetchWithAuth(downloadUrl);
      if (!imgResponse.ok) {
        throw new Error('Failed to fetch processed image');
      }
      const blob = await imgResponse.blob();
      const objectUrl = URL.createObjectURL(blob);

      const img = new Image();
      img.onload = () => {
        imageState.value = {
          ...imageState.value,
          processedWidth: img.width,
          processedHeight: img.height,
          processedSize: (blob.size / 1024).toFixed(2) + ' KB'
        };
        // Clean up the object URL after use
        URL.revokeObjectURL(objectUrl);
      };
      img.src = objectUrl;

    } catch (error) {
      imageState.value = {
        ...imageState.value,
        errorMessage: 'Failed to load processed image details'
      };
      console.error('[ERROR] Failed to load image:', error);
    }
  };

  return (
    <form onSubmit={handleSubmit} class="w-full max-w-xl mx-auto flex flex-col gap-4">
      {/* Width x Height */}
      <div class="flex flex-wrap items-center gap-2">
        <label class="font-medium text-gray-200 whitespace-nowrap">Width x Height:</label>
        <input
          type="number"
          min="1"
          class="w-20 px-2 py-1 rounded border border-gray-500 bg-gray-800 text-white"
          value={imageState.value.width}
          onInput={e => imageState.value = { ...imageState.value, width: e.currentTarget.value }}
          placeholder="Width"
        />
        <span class="text-gray-400">px x</span>
        <input
          type="number"
          min="1"
          class="w-20 px-2 py-1 rounded border border-gray-500 bg-gray-800 text-white"
          value={imageState.value.height}
          onInput={e => imageState.value = { ...imageState.value, height: e.currentTarget.value }}
          placeholder="Height"
        />
        <span class="text-gray-400">px</span>
      </div>

      {/* Maintain Aspect Ratio */}
      <div class="flex items-center gap-2">
        <input
          type="checkbox"
          id="aspect"
          checked={imageState.value.maintainAspectRatio}
          onChange={e => imageState.value = { ...imageState.value, maintainAspectRatio: e.currentTarget.checked }}
          class="accent-blue-500"
        />
        <label htmlFor="aspect" class="text-gray-200 font-medium select-none cursor-pointer">
          Maintain Aspect Ratio
        </label>
      </div>

      {/* Quality and Max Size */}
      <div class="flex flex-wrap items-center gap-2">
        <label htmlFor="quality" class="font-medium text-gray-200 whitespace-nowrap">Quality (1-100):</label>
        <input
          type="number"
          id="quality"
          min="1"
          max="100"
          class="w-20 px-2 py-1 rounded border border-gray-500 bg-gray-800 text-white"
          value={imageState.value.quality}
          onInput={e => imageState.value = { ...imageState.value, quality: Number(e.currentTarget.value) }}
          placeholder="Quality"
        />
        <span class="text-gray-400">%</span>
      </div>
      <div class="flex flex-wrap items-center gap-2">
        <label htmlFor="maxSize" class="font-medium text-gray-200 whitespace-nowrap">Max Size (KB):</label>
        <input
          type="number"
          id="maxSize"
          min="1"
          step="any"
          class="w-24 px-2 py-1 rounded border border-gray-500 bg-gray-800 text-white"
          value={imageState.value.maxSize}
          onInput={e => imageState.value = { ...imageState.value, maxSize: e.currentTarget.value }}
          placeholder="Max Size"
        />
        <span class="text-gray-400">kb</span>
      </div>
      <button
        type="submit"
        disabled={
          !imageState.value.imageFile ||
          (imageState.value.progress > 0 && imageState.value.progress < 100)
        }
        class="w-full py-2 rounded bg-blue-600 hover:bg-blue-700 text-white font-semibold transition disabled:bg-gray-500 disabled:cursor-not-allowed"
      >
        {imageState.value.progress === 100
          ? "Resize Again"
          : (imageState.value.progress > 0 && imageState.value.progress < 100)
            ? "Processing..."
            : "Resize"}
      </button>
    </form>
  );
}