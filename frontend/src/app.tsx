import { ImageUploader } from './components/ImageUploader';
import { ImageForm } from './components/ImageForm';
import { ProcessedImage } from './components/ProcessedImage';
import { imageState } from './hooks/useImageProcessor';

export function App() {
  return (
    <main class="min-h-screen bg-gray-950 font-sans p-8">
      <div class="grid grid-cols-12 gap-4 h-[90vh]">
        <div class="col-span-12 md:col-span-5 bg-gray-700 p-4 rounded-lg h-[90vh] overflow-y-auto flex flex-col gap-4">
          <div>
            <h2 class="text-3xl text-white mb-2">Single Image Resize</h2>
          </div>

          {imageState.value.errorMessage && (
            <div class="bg-red-600 text-white p-4 rounded mb-2">
              <p>{imageState.value.errorMessage}</p>
            </div>
          )}

          <div>
            <ImageUploader />
          </div>

          {imageState.value.originalWidth && (
            <div class="mt-2 p-4 border border-gray-600 bg-gray-800 rounded">
              <h3 class="text-white mb-2">Original Image Data</h3>
              <p class="text-gray-300">Width: {imageState.value.originalWidth}px</p>
              <p class="text-gray-300">Height: {imageState.value.originalHeight}px</p>
              <p class="text-gray-300">Size: {imageState.value.originalSize} KB</p>
            </div>
          )}

          <div>
            <ImageForm />
          </div>
        </div>
        <div class="col-span-12 md:col-span-7 bg-gray-900 p-4 rounded-lg flex flex-col gap-4">
          {/* Display the selected image */}
          {imageState.value.imageFile && !imageState.value.downloadUrl && (
            <div class="flex justify-center items-center w-full">
              <img
                src={URL.createObjectURL(imageState.value.imageFile)}
                alt="Selected"
                class="w-full max-h-[30vw] object-contain rounded-lg shadow-lg"
              />
            </div>
          )}
          {/* Display Processed image which we get as a presigned minio link */}
          {imageState.value.downloadUrl && (
            <div class="flex justify-center items-center w-full">
              <img
                src={imageState.value.downloadUrl}
                alt="Processed"
                class="w-full max-h-[30vw] object-contain rounded-lg shadow-lg"
              />
            </div>
          )}
          {imageState.value.progress > 0 && imageState.value.progress < 100 && (
            <div>
              <p class="text-gray-200">Progress: {imageState.value.progress}%</p>
            </div>
          )}
          {imageState.value.progress == 100 &&imageState.value.downloadUrl && (
            <div>
              <ProcessedImage />
            </div>
          )}
        </div>
      </div>
    </main>
  );
}