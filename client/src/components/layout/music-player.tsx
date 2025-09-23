import { useState, useRef, useEffect } from 'react';
import { Play, Pause, Loader2 } from 'lucide-react';

interface MusicPlayerProps {
    title: string;
    artist: string;
    albumArt: string;
    audioSrc: string;
    duration?: string;
}

/**
 * Known Issues:
 * 
 * Safari:
 * - Audio playback is unreliable — the track often stops prematurely or loops back
 *   near the beginning (around the first ~5 seconds).
 * - The `audio.duration` property frequently returns `NaN` or `0`, preventing
 *   accurate display of total track length.
 * 
 * Chrome & Firefox:
 * - Playback seems limited to the initially buffered portion of the file.
 *   After the buffered segment finishes, playback stops and resets to the "idle" state.
 * - This behavior suggests incomplete or partial loading of the audio source.
 */

export default function MusicPlayer({
    title,
    artist,
    albumArt,
    audioSrc,
}: MusicPlayerProps) {
    const [playerState, setPlayerState] = useState<'idle' | 'loading' | 'playing' | 'paused'>('idle');
    const [currentTime, setCurrentTime] = useState(0);
    const [duration, setDuration] = useState(0);
    const [isDragging, setIsDragging] = useState(false);

    const audioRef = useRef<HTMLAudioElement>(null);

    useEffect(() => {
        const audio = audioRef.current;
        if (!audio) return;

        const handleLoadedMetadata = () => {
            setDuration(audio.duration);
            setPlayerState('idle');
        };

        const handleTimeUpdate = () => {
            if (!isDragging) {
                setCurrentTime(audio.currentTime);
            }
        };

        const handleEnded = () => {
            setPlayerState('idle');
            setCurrentTime(0);
        };

        const handleCanPlay = () => {
            // Safari-specific note:
            // Avoid triggering autoplay here; Safari blocks programmatic playback
            // without a direct user gesture. Leave playback initiation to handlePlayPause().
        };

        audio.addEventListener('loadedmetadata', handleLoadedMetadata);
        audio.addEventListener('timeupdate', handleTimeUpdate);
        audio.addEventListener('ended', handleEnded);
        audio.addEventListener('canplay', handleCanPlay);

        return () => {
            audio.removeEventListener('loadedmetadata', handleLoadedMetadata);
            audio.removeEventListener('timeupdate', handleTimeUpdate);
            audio.removeEventListener('ended', handleEnded);
            audio.removeEventListener('canplay', handleCanPlay);
        };
    }, [playerState, isDragging]);

    const handlePlayPause = async () => {
        const audio = audioRef.current;
        if (!audio) return;

        try {
            if (playerState === 'idle' || playerState === 'paused') {
                if (playerState === 'idle') {
                    setPlayerState('loading');
                    // Safari compatibility: don't explicitly call audio.load(),
                    // just play directly after user interaction.
                    const playPromise = audio.play();
                    if (playPromise !== undefined) {
                        await playPromise;
                        setPlayerState('playing');
                    }
                } else {
                    setPlayerState('playing');
                    const playPromise = audio.play();
                    if (playPromise !== undefined) {
                        await playPromise;
                    }
                }
            } else if (playerState === 'playing') {
                setPlayerState('paused');
                audio.pause();
            }
        } catch (error) {
            console.error('Error playing audio:', error);
            setPlayerState('idle');
        }
    };

    const handleProgressChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const audio = audioRef.current;
        if (!audio) return;

        const newTime = (parseFloat(e.target.value) / 100) * duration;
        setCurrentTime(newTime);
        audio.currentTime = newTime;
    };

    const handleProgressMouseDown = () => {
        setIsDragging(true);
    };

    const handleProgressMouseUp = () => {
        setIsDragging(false);
    };

    const formatTime = (seconds: number) => {
        if (isNaN(seconds)) return '0:00';
        const mins = Math.floor(seconds / 60);
        const secs = Math.floor(seconds % 60);
        return `${mins}:${secs.toString().padStart(2, '0')}`;
    };

    const progressPercentage = duration > 0 ? (currentTime / duration) * 100 : 0;

    return (
        <div className="bg-gray-900 border border-gray-700 rounded-lg p-4 w-full mx-auto">
            <audio ref={audioRef} src={audioSrc} preload="none" />

            <div className="flex items-center gap-4">
                {/* Album Art */}
                <div className="flex-shrink-0">
                    <img
                        src={albumArt}
                        alt={`${title} by ${artist}`}
                        className="w-28 aspect-square rounded-md object-cover"
                    />
                </div>

                {/* Track Info and Controls */}
                <div className="flex-1 min-w-0">
                    <h3 className="text-white font-medium text-sm truncate mb-1">
                        {title}
                    </h3>
                    <p className="text-gray-400 text-sm truncate mb-3">
                        {artist}
                    </p>

                    {/* Controls and Progress */}
                    {playerState === 'idle' && (
                        <div className="flex items-center gap-3">
                            <button
                                onClick={handlePlayPause}
                                className="bg-orange-500 hover:bg-orange-400 text-black px-4 py-1.5 rounded text-sm font-medium transition-colors flex items-center gap-2"
                            >
                                <Play size={14} />
                                Play
                            </button>
                        </div>
                    )}

                    {playerState === 'loading' && (
                        <div className="flex items-center gap-2">
                            <Loader2 size={20} className="text-orange-400 animate-spin" />
                            <span className="text-gray-300 text-sm">Loading...</span>
                        </div>
                    )}

                    {(playerState === 'playing' || playerState === 'paused') && (
                        <div className="space-y-2">
                            <div className="flex items-center gap-2">
                                <button
                                    onClick={handlePlayPause}
                                    className="text-white hover:text-orange-400 transition-colors"
                                >
                                    {playerState === 'playing' ? <Pause size={20} /> : <Play size={20} />}
                                </button>
                                <span className="text-gray-300 text-xs font-mono">
                                    {formatTime(currentTime)}
                                </span>

                                {/* Progress Bar */}
                                <div className="flex-1 mx-2 relative">
                                    <div className="h-1 bg-gray-700 rounded-full overflow-hidden">
                                        <div
                                            className="h-full bg-orange-400 transition-all duration-100 ease-linear"
                                            style={{ width: `${progressPercentage}%` }}
                                        />
                                    </div>
                                    <input
                                        type="range"
                                        min="0"
                                        max="100"
                                        value={progressPercentage}
                                        onChange={handleProgressChange}
                                        onMouseDown={handleProgressMouseDown}
                                        onMouseUp={handleProgressMouseUp}
                                        className="absolute inset-0 w-full h-full opacity-0 cursor-pointer"
                                    />
                                </div>

                                <span className="text-gray-300 text-xs font-mono">
                                    -{formatTime(duration - currentTime)}
                                </span>
                            </div>
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
}
