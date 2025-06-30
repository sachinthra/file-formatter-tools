import { ImageUploader } from './components/ImageUploader';
import { ImageForm } from './components/ImageForm';
import { ProcessedImage } from './components/ProcessedImage';
import { imageState } from './hooks/useImageProcessor';
import './app.css';

export function App() {
  return (
    <main>
      <h1>Single Image Resize</h1>
      {imageState.value.errorMessage && (
        <p class="error">{imageState.value.errorMessage}</p>
      )}

      <ImageUploader />

      {imageState.value.originalWidth && (
        <div class="original-data">
          <h3>Original Image Data</h3>
          <p>Width: {imageState.value.originalWidth}px</p>
          <p>Height: {imageState.value.originalHeight}px</p>
          <p>Size: {imageState.value.originalSize}</p>
        </div>
      )}

      <ImageForm />

      {imageState.value.progress > 0 && (
        <p>Progress: {imageState.value.progress}%</p>
      )}

      {imageState.value.downloadUrl && <ProcessedImage />}
    </main>
  );
}