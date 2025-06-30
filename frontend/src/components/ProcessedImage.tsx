import { imageState } from '../hooks/useImageProcessor';
import styles from './ProcessedImage.module.css';

export function ProcessedImage() {
  return (
    <div class={styles.processedData}>
      <h3>Processed Image Data</h3>
      {imageState.value.processedWidth && (
        <>
          <p>Width: {imageState.value.processedWidth}px</p>
          <p>Height: {imageState.value.processedHeight}px</p>
          <p>Size: {imageState.value.processedSize}</p>
        </>
      )}
      <p>
        <a
          href={imageState.value.downloadUrl ?? undefined}
          target="_blank"
          rel="noopener noreferrer"
        >
          Download Resized Image
        </a>
      </p>
    </div>
  );
}