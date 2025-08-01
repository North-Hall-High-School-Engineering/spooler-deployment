import { useEffect, useMemo, useRef } from "react";
import { Canvas, useThree } from "@react-three/fiber";
import { PerspectiveCamera } from "@react-three/drei";
import { STLLoader } from "three-stdlib";
import * as THREE from "three";

type SimpleStlPreviewProps = {
  stlData: string; // base64 string
  color?: string;
  width?: number;
  height?: number;
};

const StlViewer = ({
  stlData,
  color,
}: {
  stlData: string;
  color: string;
}) => {
  const meshRef = useRef<THREE.Mesh>(null);
  const cameraRef = useRef<THREE.PerspectiveCamera>(null);
  const { size } = useThree();

  const geometry = useMemo(() => {
    const loader = new STLLoader();
    const binary = atob(stlData);
    const buffer = new Uint8Array(binary.length);
    for (let i = 0; i < binary.length; i++) {
      buffer[i] = binary.charCodeAt(i);
    }
    const geom = loader.parse(buffer.buffer);
    geom.computeVertexNormals();
    geom.computeBoundingBox();
    return geom;
  }, [stlData]);

  useEffect(() => {
    if (!cameraRef.current || !geometry.boundingBox) return;

    const box = geometry.boundingBox;
    const center = new THREE.Vector3();
    const sizeVec = new THREE.Vector3();
    box.getCenter(center);
    box.getSize(sizeVec);

    // Re-center geometry
    geometry.translate(-center.x, -center.y, -center.z);

    // Position camera
    const maxDim = Math.max(sizeVec.x, sizeVec.y, sizeVec.z);
    const distance = maxDim / (2 * Math.tan((cameraRef.current.fov * Math.PI) / 360));

    const cameraPos = new THREE.Vector3(-distance, distance, distance);
    cameraRef.current.position.copy(cameraPos);
    cameraRef.current.near = 0.01;
    cameraRef.current.far = distance * 10;
    cameraRef.current.lookAt(new THREE.Vector3(0, 0, 0));
    cameraRef.current.updateProjectionMatrix();
  }, [geometry]);

  return (
    <>
      <PerspectiveCamera
        ref={cameraRef}
        makeDefault
        fov={45}
        aspect={size.width / size.height}
      />
      <ambientLight intensity={0.6} />
      <directionalLight position={[5, 5, 5]} intensity={0.9} />
      <mesh ref={meshRef} geometry={geometry}>
        <meshStandardMaterial color={color} />
      </mesh>
    </>
  );
};

export const SimpleStlPreview = ({
  stlData,
  color = "#000000",
  width = 300,
  height = 250,
}: SimpleStlPreviewProps) => (
  <Canvas style={{ width, height }}>
    <StlViewer stlData={stlData} color={color} />
  </Canvas>
);
